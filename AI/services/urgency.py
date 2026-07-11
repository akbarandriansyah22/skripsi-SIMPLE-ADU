"""Penentuan tingkat urgensi berbasis aturan."""

from __future__ import annotations

import logging

logger = logging.getLogger(__name__)

HIGH_URGENCY_THRESHOLD = 5
MEDIUM_URGENCY_THRESHOLD = 2
NEGATIVE_SCORE_HIGH_THRESHOLD = -4
NEGATIVE_SCORE_MEDIUM_THRESHOLD = -2

HIGH_URGENCY_KEYWORDS = {
    "bahaya",
    "darurat",
    "fatal",
    "gawat",
    "kecelakaan",
    "kebakaran",
    "krisis",
    "parah",
    "penting",
    "rusak",
    "segera",
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
    urgency_score = calculate_urgency_score(tokens, sentiment_score)

    if urgency_score >= HIGH_URGENCY_THRESHOLD:
        return "Tinggi"
    if urgency_score >= MEDIUM_URGENCY_THRESHOLD:
        return "Sedang"
    return "Rendah"
