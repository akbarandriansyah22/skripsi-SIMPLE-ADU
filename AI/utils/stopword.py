"""Stopword removal Bahasa Indonesia menggunakan NLTK."""

from __future__ import annotations

import logging
from functools import lru_cache

from nltk.corpus import stopwords

logger = logging.getLogger(__name__)

FALLBACK_INDONESIAN_STOPWORDS = {
    "ada",
    "adalah",
    "agar",
    "akan",
    "aku",
    "anda",
    "atau",
    "bagai",
    "bagaimana",
    "bagi",
    "bahwa",
    "baik",
    "banyak",
    "baru",
    "bisa",
    "dan",
    "dari",
    "dengan",
    "di",
    "dia",
    "itu",
    "jadi",
    "jika",
    "juga",
    "karena",
    "ke",
    "kita",
    "lagi",
    "lebih",
    "maka",
    "mereka",
    "namun",
    "pada",
    "saat",
    "saja",
    "saya",
    "sebagai",
    "seperti",
    "serta",
    "sudah",
    "supaya",
    "tapi",
    "telah",
    "tentang",
    "tersebut",
    "tidak",
    "untuk",
    "yang",
}


@lru_cache(maxsize=1)
def get_indonesian_stopwords() -> set[str]:
    """Ambil stopword Bahasa Indonesia dari NLTK dengan fallback lokal."""
    try:
        return set(stopwords.words("indonesian"))
    except LookupError:
        logger.warning("Resource stopwords NLTK belum tersedia, memakai fallback lokal")
        return FALLBACK_INDONESIAN_STOPWORDS
    except Exception:
        logger.exception("Gagal memuat stopwords NLTK, memakai fallback lokal")
        return FALLBACK_INDONESIAN_STOPWORDS


def remove_stopwords(tokens: list[str]) -> list[str]:
    """Hapus stopword dari daftar token."""
    stopword_set = get_indonesian_stopwords()
    return [token for token in tokens if token and token not in stopword_set]
