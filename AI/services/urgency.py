"""Penentuan tingkat urgensi berbasis aturan."""

from __future__ import annotations

import logging

from services.urgency_model import URGENCY_LEVELS, predict_urgency_from_model
from services.negation import effective_tokens

logger = logging.getLogger(__name__)

HIGH_URGENCY_THRESHOLD = 5
MEDIUM_URGENCY_THRESHOLD = 2
NEGATIVE_SCORE_HIGH_THRESHOLD = -4
NEGATIVE_SCORE_MEDIUM_THRESHOLD = -2

HIGH_URGENCY_KEYWORDS = {
    "ancam",
    "bahaya",
    "bakar",
    "darurat",
    "fatal",
    "gawat",
    "jatuh",
    "kecelakaan",
    "kebakaran",
    "korsleting",
    "krisis",
    "luka",
    "parah",
    "penting",
    "rusak",
    "segera",
    "terluka",
}

CRITICAL_URGENCY_KEYWORDS = {
    "api",
    "asap",
    "bakar",
    "darurat",
    "gawat",
    "kebakaran",
    "korsleting",
    "luka",
    "pingsan",
    "terbakar",
    "terluka",
}

MEDIUM_URGENCY_KEYWORDS = {
    "ganggu",
    "hambat",
    "kendala",
    "keluhan",
    "lambat",
    "masalah",
    "mengganggu",
    "respon",
    "sulit",
    "terkendala",
}


def calculate_urgency_score(tokens: list[str], sentiment_score: int = 0) -> int:
    """Hitung skor urgensi berdasarkan keyword dan skor sentimen."""
    token_set = set(tokens)
    keyword_score = 0
    keyword_score += sum(2 for keyword in HIGH_URGENCY_KEYWORDS if keyword in token_set)
    keyword_score += sum(1 for keyword in MEDIUM_URGENCY_KEYWORDS if keyword in token_set)

    if sentiment_score <= NEGATIVE_SCORE_HIGH_THRESHOLD:
        keyword_score += 2
    elif sentiment_score <= NEGATIVE_SCORE_MEDIUM_THRESHOLD:
        keyword_score += 1

    logger.info("Skor urgensi: %s", keyword_score)
    return keyword_score


def determine_urgency(tokens: list[str], sentiment_score: int = 0) -> str:
    """Tentukan tingkat urgensi Rendah, Sedang, atau Tinggi."""
    active_tokens = effective_tokens(tokens)
    token_set = set(active_tokens)
    if token_set.intersection(CRITICAL_URGENCY_KEYWORDS):
        logger.info("Urgensi tinggi karena keyword kritis terdeteksi")
        return "Tinggi"

    urgency_score = calculate_urgency_score(active_tokens, sentiment_score)

    rule_urgency = "Rendah"
    if urgency_score >= HIGH_URGENCY_THRESHOLD:
        rule_urgency = "Tinggi"
    elif urgency_score >= MEDIUM_URGENCY_THRESHOLD:
        rule_urgency = "Sedang"

    model_urgency = predict_urgency_from_model(active_tokens)
    if model_urgency in URGENCY_LEVELS:
        return max((rule_urgency, model_urgency), key=lambda label: URGENCY_LEVELS[label])

    return rule_urgency
