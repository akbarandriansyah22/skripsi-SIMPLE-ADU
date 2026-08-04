"""Analisis sentimen berbasis InSet Lexicon."""

from __future__ import annotations

import csv
import logging
from functools import lru_cache
from typing import Any
from pathlib import Path

from services.negation import is_negated

logger = logging.getLogger(__name__)

LEXICON_DIR = Path(__file__).resolve().parents[1] / "lexicon"
POSITIVE_LEXICON_PATH = LEXICON_DIR / "positive.tsv"
NEGATIVE_LEXICON_PATH = LEXICON_DIR / "negative.tsv"


def _parse_score(raw_score: str, default: int) -> int:
    """Ubah skor TSV menjadi integer."""
    try:
        return int(float(raw_score))
    except (TypeError, ValueError):
        return default


def _load_lexicon(path: Path, default_score: int) -> dict[str, int]:
    """Baca file TSV lexicon menjadi mapping kata dan skor."""
    if not path.exists():
        logger.error("File lexicon tidak ditemukan: %s", path)
        raise FileNotFoundError(f"File lexicon tidak ditemukan: {path}")

    lexicon: dict[str, int] = {}
    try:
        with path.open("r", encoding="utf-8") as file:
            reader = csv.reader(file, delimiter="\t")
            for row in reader:
                if not row:
                    continue
                word = row[0].strip().lower()
                if not word or word in {"word", "kata"}:
                    continue
                score = _parse_score(row[1], default_score) if len(row) > 1 else default_score
                lexicon[word] = score
    except Exception:
        logger.exception("Gagal membaca lexicon: %s", path)
        raise

    logger.info("Berhasil memuat %s kata dari %s", len(lexicon), path.name)
    return lexicon


@lru_cache(maxsize=1)
def load_sentiment_lexicons() -> tuple[dict[str, int], dict[str, int]]:
    """Muat positive dan negative lexicon."""
    positive_words = _load_lexicon(POSITIVE_LEXICON_PATH, 1)
    negative_words = _load_lexicon(NEGATIVE_LEXICON_PATH, -1)
    return positive_words, negative_words


def analyze_sentiment(tokens: list[str]) -> dict[str, Any]:
    """Hitung skor sentimen dari daftar token."""
    try:
        positive_words, negative_words = load_sentiment_lexicons()
    except FileNotFoundError:
        raise
    except Exception as exc:
        logger.exception("Analisis sentimen gagal saat memuat lexicon")
        raise ValueError(f"Analisis sentimen gagal: {exc}") from exc

    positive_score = 0
    negative_score = 0
    matched_words: list[dict[str, Any]] = []
    for index, token in enumerate(tokens):
        positive = positive_words.get(token, 0)
        negative = negative_words.get(token, 0)
        if is_negated(tokens, index):
            positive, negative = -negative, -positive
        positive_score += positive
        negative_score += negative
        if positive or negative:
            matched_words.append({
                "word": token,
                "positive_score": positive,
                "negative_score": negative,
                "negated": is_negated(tokens, index),
            })
    total_score = positive_score + negative_score

    if total_score > 0:
        sentiment = "Positif"
    elif total_score < 0:
        sentiment = "Negatif"
    else:
        sentiment = "Netral"

    return {
        "positive_score": positive_score,
        "negative_score": negative_score,
        "score": total_score,
        "sentiment": sentiment,
        "matched_words": matched_words,
        "explanation": (
            f"Ditemukan {len(matched_words)} kata InSet yang memengaruhi skor; "
            f"skor positif {positive_score}, skor negatif {negative_score}, "
            f"skor akhir {total_score}."
        ),
    }
