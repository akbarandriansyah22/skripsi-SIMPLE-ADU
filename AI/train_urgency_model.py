"""Train and evaluate the Multinomial Naive Bayes urgency model."""

from __future__ import annotations

import json
import math
import random
from collections import Counter, defaultdict
from datetime import datetime, timezone
from pathlib import Path
from typing import Iterable

from services.preprocessing import preprocess_text
from utils.read_excel import read_dataset

REVIEWED_DATASET_PATH = Path(__file__).resolve().parent / "data" / "dataset_urgency_reviewed.xlsx"
MODEL_PATH = Path(__file__).resolve().parent / "models" / "urgency_model.json"
AUDIT_REPORT_PATH = Path(__file__).resolve().parent / "reports" / "urgency_dataset_audit.md"
TEXT_COLUMN = "teks_pengaduan"
LABEL_COLUMN = "label_reviewed"
RANDOM_SEED = 42
LABELS = ("Rendah", "Sedang", "Tinggi")


def _tokenized_rows(dataframe) -> tuple[list[tuple[list[str], str]], int]:
    rows: list[tuple[list[str], str]] = []
    skipped = 0
    for _, row in dataframe.iterrows():
        text = row.get(TEXT_COLUMN)
        label = str(row.get(LABEL_COLUMN, "")).strip()
        if text is None or not str(text).strip() or label not in LABELS:
            skipped += 1
            continue
        tokens = preprocess_text(str(text))["tokens"]
        if not isinstance(tokens, list) or not tokens:
            skipped += 1
            continue
        rows.append((tokens, label))
    return rows, skipped


def _fit_nb(rows: Iterable[tuple[list[str], str]]) -> dict:
    class_doc_counts: Counter[str] = Counter()
    class_token_counts: dict[str, Counter[str]] = defaultdict(Counter)
    vocabulary: set[str] = set()

    for tokens, label in rows:
        class_doc_counts[label] += 1
        class_token_counts[label].update(tokens)
        vocabulary.update(tokens)

    if not class_doc_counts or not vocabulary:
        raise ValueError("Tidak ada data latih yang valid")

    classes = sorted(class_doc_counts)
    total_docs = sum(class_doc_counts.values())
    vocab_size = len(vocabulary)
    class_log_prior = {
        label: math.log(class_doc_counts[label] / total_docs) for label in classes
    }
    feature_log_prob: dict[str, dict[str, float]] = {}
    unknown_log_prob: dict[str, float] = {}

    for label in classes:
        total_tokens = sum(class_token_counts[label].values())
        denominator = total_tokens + vocab_size
        feature_log_prob[label] = {
            token: math.log((class_token_counts[label][token] + 1) / denominator)
            for token in sorted(vocabulary)
        }
        unknown_log_prob[label] = math.log(1 / denominator)

    return {
        "classes": classes,
        "class_counts": dict(class_doc_counts),
        "vocabulary_size": vocab_size,
        "class_log_prior": class_log_prior,
        "feature_log_prob": feature_log_prob,
        "unknown_log_prob": unknown_log_prob,
    }


def _predict(model: dict, tokens: list[str]) -> str:
    scores: dict[str, float] = {}
    for label in model["classes"]:
        score = float(model["class_log_prior"][label])
        unknown = float(model["unknown_log_prob"][label])
        probabilities = model["feature_log_prob"][label]
        score += sum(float(probabilities.get(token, unknown)) for token in tokens)
        scores[label] = score
    return max(scores, key=scores.get)


def _stratified_folds(rows: list[tuple[list[str], str]], seed: int) -> list[list[int]]:
    by_class: dict[str, list[int]] = defaultdict(list)
    for index, (_, label) in enumerate(rows):
        by_class[label].append(index)
    minimum_class_size = min(len(indices) for indices in by_class.values())
    fold_count = min(5, minimum_class_size)
    if fold_count < 2:
        return []

    rng = random.Random(seed)
    folds: list[list[int]] = [[] for _ in range(fold_count)]
    for indices in by_class.values():
        shuffled = list(indices)
        rng.shuffle(shuffled)
        for offset, index in enumerate(shuffled):
            folds[offset % fold_count].append(index)
    return folds


def evaluate(rows: list[tuple[list[str], str]], seed: int = RANDOM_SEED) -> dict:
    """Run deterministic stratified cross-validation without train-set scoring."""
    folds = _stratified_folds(rows, seed)
    if not folds:
        return {"method": "not_available", "reason": "Jumlah data per kelas tidak mencukupi."}

    confusion = {actual: {predicted: 0 for predicted in LABELS} for actual in LABELS}
    for validation_fold in folds:
        validation_indices = set(validation_fold)
        train_rows = [row for index, row in enumerate(rows) if index not in validation_indices]
        model = _fit_nb(train_rows)
        for index in validation_fold:
            tokens, actual = rows[index]
            predicted = _predict(model, tokens)
            confusion.setdefault(actual, {label: 0 for label in LABELS})[predicted] += 1

    total = sum(sum(row.values()) for row in confusion.values())
    correct = sum(confusion[label][label] for label in LABELS)
    metrics: dict[str, dict[str, float]] = {}
    for label in LABELS:
        true_positive = confusion[label][label]
        false_positive = sum(confusion[actual][label] for actual in LABELS if actual != label)
        false_negative = sum(confusion[label][predicted] for predicted in LABELS if predicted != label)
        precision = true_positive / (true_positive + false_positive) if true_positive + false_positive else 0.0
        recall = true_positive / (true_positive + false_negative) if true_positive + false_negative else 0.0
        f1 = 2 * precision * recall / (precision + recall) if precision + recall else 0.0
        metrics[label] = {"precision": precision, "recall": recall, "f1": f1}

    return {
        "method": "stratified_k_fold",
        "folds": len(folds),
        "random_seed": seed,
        "accuracy": correct / total if total else 0.0,
        "precision_recall_f1_per_class": metrics,
        "confusion_matrix": confusion,
    }


def train() -> dict:
    if not REVIEWED_DATASET_PATH.exists():
        from review_urgency_dataset import review_dataset

        review_dataset()

    dataframe = read_dataset(REVIEWED_DATASET_PATH)
    required_columns = {TEXT_COLUMN, LABEL_COLUMN}
    missing_columns = required_columns.difference(dataframe.columns)
    if missing_columns:
        raise ValueError(f"Kolom dataset reviewed tidak ditemukan: {sorted(missing_columns)}")

    rows, skipped_rows = _tokenized_rows(dataframe)
    model = _fit_nb(rows)
    evaluation = evaluate(rows, RANDOM_SEED)
    model.update(
        {
            "model_type": "multinomial_naive_bayes",
            "source_dataset": str(REVIEWED_DATASET_PATH.relative_to(REVIEWED_DATASET_PATH.parents[1])),
            "total_rows": len(dataframe),
            "valid_rows": len(rows),
            "skipped_rows": skipped_rows,
            "label_column": LABEL_COLUMN,
            "text_column": TEXT_COLUMN,
            "training_timestamp": datetime.now(timezone.utc).isoformat(),
            "random_seed": RANDOM_SEED,
            "curated_reference_count": 0,
            "evaluation": evaluation,
        }
    )

    MODEL_PATH.parent.mkdir(parents=True, exist_ok=True)
    MODEL_PATH.write_text(json.dumps(model, ensure_ascii=False, indent=2), encoding="utf-8")
    _append_evaluation_to_report(model)
    return model


def _append_evaluation_to_report(model: dict) -> None:
    """Keep the audit report self-contained after the model is evaluated."""
    if not AUDIT_REPORT_PATH.exists():
        return
    report = AUDIT_REPORT_PATH.read_text(encoding="utf-8")
    marker = "\n## Evaluasi model\n"
    report = report.split(marker, 1)[0]
    evaluation = model["evaluation"]
    if evaluation.get("method") == "not_available":
        section = (
            marker
            + "\nEvaluasi stratified cross-validation tidak tersedia karena jumlah data per kelas tidak mencukupi.\n"
        )
    else:
        lines = [
            marker,
            "",
            f"Metode: `{evaluation['method']}` dengan {evaluation['folds']} fold dan random seed `{evaluation['random_seed']}`.",
            f"Accuracy: `{evaluation['accuracy']:.4f}`.",
            "",
            "| Kelas | Precision | Recall | F1-score |",
            "|---|---:|---:|---:|",
        ]
        for label, metrics in evaluation["precision_recall_f1_per_class"].items():
            lines.append(
                f"| {label} | {metrics['precision']:.4f} | {metrics['recall']:.4f} | {metrics['f1']:.4f} |"
            )
        lines.extend(
            [
                "",
                "Confusion matrix:",
                "",
                "```json",
                json.dumps(evaluation["confusion_matrix"], ensure_ascii=False, indent=2),
                "```",
                "",
                "Kelas Tinggi tidak memiliki sampel reviewed, sehingga precision/recall/F1 Tinggi bernilai 0 dan tidak boleh ditafsirkan sebagai evaluasi kemampuan deteksi Tinggi. Guardrail rule-based dievaluasi melalui test kasus keselamatan dan dampak akademik kritis.",
                "",
            ]
        )
        section = "\n".join(lines)
    AUDIT_REPORT_PATH.write_text(report.rstrip() + "\n" + section, encoding="utf-8")


if __name__ == "__main__":
    trained_model = train()
    print(f"Model tersimpan: {MODEL_PATH}")
    print(f"Jumlah baris: {trained_model['total_rows']}")
    print(f"Baris valid: {trained_model['valid_rows']}")
    print(f"Baris dilewati: {trained_model['skipped_rows']}")
    print(f"Jumlah kelas: {trained_model['class_counts']}")
    print(f"Ukuran vocabulary: {trained_model['vocabulary_size']}")
    print(f"Evaluasi: {json.dumps(trained_model['evaluation'], ensure_ascii=False)}")
