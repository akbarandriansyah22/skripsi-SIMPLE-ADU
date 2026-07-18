"""Latih model urgensi dari dataset pengaduan."""

from __future__ import annotations

import json
import math
from collections import Counter, defaultdict
from pathlib import Path

from services.preprocessing import preprocess_text
from utils.read_excel import DEFAULT_DATASET_PATH, read_dataset

TEXT_COLUMN = "Tuliskan pengaduan, keluhan, atau saran yang pernah Anda alami selama menjadi mahasiswa FT UMJ.\n"
CONDITION_COLUMN = "Bagaimana kondisi yang Anda alami terkait pengaduan tersebut?"
MODEL_PATH = Path(__file__).resolve().parent / "models" / "urgency_model.json"

CONDITION_TO_URGENCY = {
    "Biasa saja": "Rendah",
    "Tidak terlalu bermasalah": "Rendah",
    "Cukup mengganggu": "Sedang",
    "Sangat merugikan": "Tinggi",
}


def train() -> dict:
    dataframe = read_dataset(DEFAULT_DATASET_PATH)
    required_columns = {TEXT_COLUMN, CONDITION_COLUMN}
    missing_columns = required_columns.difference(dataframe.columns)
    if missing_columns:
        raise ValueError(f"Kolom dataset tidak ditemukan: {sorted(missing_columns)}")

    class_doc_counts: Counter[str] = Counter()
    class_token_counts: dict[str, Counter[str]] = defaultdict(Counter)
    vocabulary: set[str] = set()

    for _, row in dataframe.iterrows():
        text = str(row[TEXT_COLUMN]).strip()
        condition = str(row[CONDITION_COLUMN]).strip()
        label = CONDITION_TO_URGENCY.get(condition)
        if not text or not label:
            continue

        tokens = preprocess_text(text)["tokens"]
        if not isinstance(tokens, list) or not tokens:
            continue

        class_doc_counts[label] += 1
        class_token_counts[label].update(tokens)
        vocabulary.update(tokens)

    if not class_doc_counts:
        raise ValueError("Tidak ada data latih yang valid")

    classes = sorted(class_doc_counts.keys())
    total_docs = sum(class_doc_counts.values())
    vocab_size = len(vocabulary)
    class_log_prior = {
        label: math.log(class_doc_counts[label] / total_docs)
        for label in classes
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

    model = {
        "model_type": "multinomial_naive_bayes",
        "text_column": TEXT_COLUMN,
        "label_column": CONDITION_COLUMN,
        "label_mapping": CONDITION_TO_URGENCY,
        "classes": classes,
        "class_counts": dict(class_doc_counts),
        "vocabulary_size": vocab_size,
        "class_log_prior": class_log_prior,
        "feature_log_prob": feature_log_prob,
        "unknown_log_prob": unknown_log_prob,
    }

    MODEL_PATH.parent.mkdir(parents=True, exist_ok=True)
    with MODEL_PATH.open("w", encoding="utf-8") as file:
        json.dump(model, file, ensure_ascii=False, indent=2)

    return model


if __name__ == "__main__":
    trained_model = train()
    print(f"Model tersimpan: {MODEL_PATH}")
    print(f"Jumlah kelas: {trained_model['class_counts']}")
    print(f"Ukuran vocabulary: {trained_model['vocabulary_size']}")
