"""Model urgensi sederhana yang dilatih dari dataset pengaduan."""

from __future__ import annotations

import json
import logging
import math
from functools import lru_cache
from pathlib import Path

logger = logging.getLogger(__name__)

MODEL_PATH = Path(__file__).resolve().parents[1] / "models" / "urgency_model.json"
URGENCY_LEVELS = {"Rendah": 0, "Sedang": 1, "Tinggi": 2}


@lru_cache(maxsize=1)
def load_urgency_model() -> dict | None:
    """Muat artifact model urgency jika sudah tersedia."""
    if not MODEL_PATH.exists():
        logger.warning("Model urgensi belum tersedia: %s", MODEL_PATH)
        return None

    try:
        with MODEL_PATH.open("r", encoding="utf-8") as file:
            return json.load(file)
    except Exception:
        logger.exception("Gagal memuat model urgensi")
        return None


def predict_urgency_from_model(tokens: list[str]) -> str | None:
    """Prediksi urgency dengan Multinomial Naive Bayes sederhana."""
    model = load_urgency_model()
    if not model:
        return None

    classes = model.get("classes", [])
    class_log_prior = model.get("class_log_prior", {})
    feature_log_prob = model.get("feature_log_prob", {})
    unknown_log_prob = model.get("unknown_log_prob", {})

    if not classes or not class_log_prior or not feature_log_prob:
        logger.warning("Artifact model urgensi tidak lengkap")
        return None

    scores: dict[str, float] = {}
    for label in classes:
        score = float(class_log_prior.get(label, -math.inf))
        token_probs = feature_log_prob.get(label, {})
        unknown_prob = float(unknown_log_prob.get(label, -20.0))
        for token in tokens:
            score += float(token_probs.get(token, unknown_prob))
        scores[label] = score

    if not scores:
        return None

    return max(scores, key=scores.get)
