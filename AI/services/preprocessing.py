"""Pipeline preprocessing teks pengaduan."""

from __future__ import annotations

import logging

from Sastrawi.Stemmer.StemmerFactory import StemmerFactory

from utils.stopword import remove_stopwords
from utils.text_cleaning import clean_text

logger = logging.getLogger(__name__)

_stemmer = StemmerFactory().create_stemmer()


def tokenize(text: str) -> list[str]:
    """Pisahkan teks menjadi token berbasis spasi."""
    return text.split()


def stem_tokens(tokens: list[str]) -> list[str]:
    """Lakukan stemming PySastrawi pada setiap token."""
    stemmed_tokens: list[str] = []
    for token in tokens:
        try:
            stemmed_tokens.append(_stemmer.stem(token))
        except Exception:
            logger.exception("Gagal melakukan stemming untuk token: %s", token)
            stemmed_tokens.append(token)
    return stemmed_tokens


def preprocess_text(text: str) -> dict[str, str | list[str]]:
    """Jalankan pipeline cleaning, tokenizing, stopword removal, dan stemming."""
    if not isinstance(text, str):
        logger.warning("Input preprocessing bukan string: %s", type(text).__name__)
        text = "" if text is None else str(text)

    cleaned = clean_text(text)
    tokens = tokenize(cleaned)
    filtered_tokens = remove_stopwords(tokens)
    stemmed_tokens = stem_tokens(filtered_tokens)
    cleaned_text = " ".join(stemmed_tokens)

    return {
        "cleaned_text": cleaned_text,
        "tokens": stemmed_tokens,
    }
