"""Fungsi pembersihan teks untuk preprocessing Bahasa Indonesia."""

from __future__ import annotations

import logging
import re
import string

logger = logging.getLogger(__name__)

URL_PATTERN = re.compile(r"https?://\S+|www\.\S+")
EMOJI_PATTERN = re.compile(
    "["
    "\U0001f600-\U0001f64f"
    "\U0001f300-\U0001f5ff"
    "\U0001f680-\U0001f6ff"
    "\U0001f1e0-\U0001f1ff"
    "\U00002700-\U000027bf"
    "\U0001f900-\U0001f9ff"
    "\U00002600-\U000026ff"
    "]+",
    flags=re.UNICODE,
)
PUNCTUATION_TRANSLATOR = str.maketrans("", "", string.punctuation)


def case_folding(text: str) -> str:
    """Ubah teks menjadi huruf kecil."""
    return text.lower()


def remove_url(text: str) -> str:
    """Hapus URL dari teks."""
    return URL_PATTERN.sub(" ", text)


def remove_numbers(text: str) -> str:
    """Hapus angka dari teks."""
    return re.sub(r"\d+", " ", text)


def remove_punctuation(text: str) -> str:
    """Hapus tanda baca dari teks."""
    return text.translate(PUNCTUATION_TRANSLATOR)


def remove_emoji(text: str) -> str:
    """Hapus emoji dan simbol visual umum dari teks."""
    return EMOJI_PATTERN.sub(" ", text)


def normalize_whitespace(text: str) -> str:
    """Rapikan spasi berlebih."""
    return re.sub(r"\s+", " ", text).strip()


def clean_text(text: str) -> str:
    """Jalankan seluruh tahapan text cleaning."""
    if not isinstance(text, str):
        logger.warning("Input text cleaning bukan string: %s", type(text).__name__)
        text = "" if text is None else str(text)

    cleaned = case_folding(text)
    cleaned = remove_url(cleaned)
    cleaned = remove_emoji(cleaned)
    cleaned = remove_numbers(cleaned)
    cleaned = remove_punctuation(cleaned)
    cleaned = normalize_whitespace(cleaned)
    return cleaned
