"""Service utama untuk prediksi sentimen dan urgensi."""

from __future__ import annotations

import logging
from typing import Any

from services.preprocessing import preprocess_text
from services.sentiment import analyze_sentiment
from services.urgency import urgency_analysis

logger = logging.getLogger(__name__)


def predict_complaint(deskripsi: str) -> dict[str, Any]:
    """Gabungkan preprocessing, sentiment analysis, dan urgency analysis."""
    if not isinstance(deskripsi, str) or not deskripsi.strip():
        logger.warning("Deskripsi kosong atau tidak valid")
        raise ValueError("Deskripsi wajib diisi")

    preprocessing_result = preprocess_text(deskripsi)
    tokens = preprocessing_result["tokens"]

    if not isinstance(tokens, list):
        logger.error("Hasil token preprocessing tidak valid")
        raise ValueError("Hasil preprocessing tidak valid")

    sentiment_result = analyze_sentiment(tokens)
    score = int(sentiment_result["score"])
    urgency_result = urgency_analysis(tokens, score)

    return {
        "original_text": deskripsi,
        "cleaned_text": preprocessing_result["cleaned_text"],
        "tokens": tokens,
        "score": score,
        "sentiment": sentiment_result["sentiment"],
        "positive_score": sentiment_result["positive_score"],
        "negative_score": sentiment_result["negative_score"],
        "sentiment_score": score,
        "sentiment_label": sentiment_result["sentiment"],
        "matched_words": sentiment_result["matched_words"],
        "sentiment_explanation": sentiment_result["explanation"],
        "urgency_score": urgency_result["score"],
        "urgency_label": urgency_result["label"],
        "urgency_reason": urgency_result["reason"],
        "urgency": urgency_result["label"],
    }
