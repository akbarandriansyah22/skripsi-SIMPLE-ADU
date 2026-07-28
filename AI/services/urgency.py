"""Penentuan urgensi berbasis konteks, guardrail, dan model MNB pendukung."""

from __future__ import annotations

import logging

from services.negation import effective_tokens, is_negated
from services.urgency_model import URGENCY_LEVELS, predict_urgency_from_model

logger = logging.getLogger(__name__)

# These groups intentionally do not reuse the sentiment lexicon as urgency
# labels.  A negative word describes polarity; it does not prove severity.
CRITICAL_SAFETY_KEYWORDS = frozenset(
    {
        "api",
        "asap",
        "bakar",
        "bahaya",
        "darurat",
        "gawat",
        "kecelakaan",
        "kebakaran",
        "korsleting",
        "luka",
        "pingsan",
        "roboh",
        "terbakar",
        "terluka",
    }
)

# Tuples are used as token patterns after the normal project preprocessing
# (cleaning, stopword removal, and Sastrawi stemming).
HIGH_ACADEMIC_IMPACT_PATTERNS = {
    "cakupan_luas": frozenset({"seluruh", "banyak", "mahasiswa"}),
    "batas_waktu_kritis": frozenset({"batas", "akhir"}),
    "tidak_ada_alternatif": frozenset({"alternatif"}),
    "akses_total_terhenti": frozenset({"akses", "berhenti", "total", "sama", "sekali"}),
    "kehilangan_hak": frozenset({"gagal", "hak", "ambil", "kuliah"}),
}

MEDIUM_OPERATIONAL_KEYWORDS = frozenset(
    {
        "berisik",
        "down",
        "gagal",
        "ganggu",
        "hambat",
        "kendala",
        "kotor",
        "lama",
        "lambat",
        "login",
        "mati",
        "masalah",
        "mengganggu",
        "panas",
        "rusak",
        "sulit",
        "telat",
        "terbatas",
        "terhambat",
        "terlambat",
        "terkendala",
    }
)

MITIGATING_PHRASES = (
    ("sedikit",),
    ("sesekali",),
    ("sementara",),
    ("masih", "dapat", "guna"),
    ("tetap", "jalan"),
    ("tidak", "terlalu", "ganggu"),
    ("ada", "alternatif"),
    ("hanya", "saran"),
    ("sekadar", "tanya"),
)

NEGATION_TERMS = frozenset({"tidak", "bukan", "tak", "belum", "jangan", "tanpa", "tiada", "kurang"})

# Kept for callers that used the previous helper.  The numeric value is only
# an operational-signal summary; it is never used as a direct high-urgency
# threshold and sentiment contributes at most one bounded point.
NEGATIVE_SENTIMENT_SUPPORT = -1


def _contains_phrase(tokens: list[str], phrase: tuple[str, ...]) -> bool:
    width = len(phrase)
    return any(tuple(tokens[index : index + width]) == phrase for index in range(len(tokens) - width + 1))


def _negated_token_present(tokens: list[str], candidates: set[str] | frozenset[str]) -> bool:
    return any(token in candidates and is_negated(tokens, index) for index, token in enumerate(tokens))


def analyze_urgency_features(tokens: list[str]) -> dict[str, object]:
    """Extract explainable urgency signals while respecting short negations."""
    active_tokens = effective_tokens(tokens)
    active_set = set(active_tokens)
    critical = sorted(active_set.intersection(CRITICAL_SAFETY_KEYWORDS))
    medium = sorted(active_set.intersection(MEDIUM_OPERATIONAL_KEYWORDS))

    broad_scope = bool(active_set.intersection(HIGH_ACADEMIC_IMPACT_PATTERNS["cakupan_luas"]))
    critical_deadline = HIGH_ACADEMIC_IMPACT_PATTERNS["batas_waktu_kritis"].issubset(active_set)
    no_alternative = _negated_token_present(tokens, HIGH_ACADEMIC_IMPACT_PATTERNS["tidak_ada_alternatif"])
    total_outage = (
        ("akses" in tokens and {"sama", "sekali"}.issubset(set(tokens)))
        or {"berhenti", "total"}.issubset(set(tokens))
    )
    academic_loss = (
        {"gagal", "ambil"}.issubset(set(tokens))
        or {"kehilangan", "hak"}.issubset(set(tokens))
        or "hak" in set(tokens) and "akademik" in set(tokens)
    )
    high_academic_impact = (
        broad_scope
        and critical_deadline
        and no_alternative
        and (total_outage or academic_loss)
    )

    operational_block = bool(
        _negated_token_present(
            tokens,
            {"akses", "fungsi", "guna", "login", "pakai", "jalan", "berfungsi"},
        )
    )
    mitigating = [
        " ".join(phrase)
        for phrase in MITIGATING_PHRASES
        if _contains_phrase(tokens, phrase)
    ]

    return {
        "active_tokens": active_tokens,
        "critical_keywords": critical,
        "medium_keywords": medium,
        "high_impact": high_academic_impact,
        "high_impact_indicators": {
            "cakupan_luas": broad_scope,
            "batas_waktu_kritis": critical_deadline,
            "tidak_ada_alternatif": no_alternative,
            "layanan_berhenti_total": total_outage,
            "kehilangan_hak_akademik": academic_loss,
        },
        "operational_block": operational_block,
        "mitigating_phrases": mitigating,
        "negation_detected": any(token in NEGATION_TERMS for token in tokens),
    }


def calculate_urgency_score(tokens: list[str], sentiment_score: int = 0) -> int:
    """Return a bounded operational support score for diagnostics only."""
    features = analyze_urgency_features(tokens)
    score = len(features["medium_keywords"]) + int(features["operational_block"])
    if sentiment_score < 0:
        score += 1
    logger.info("Skor sinyal urgensi operasional: %s", score)
    return score


def determine_urgency(tokens: list[str], sentiment_score: int = 0) -> str:
    """Determine urgency independently from sentiment polarity."""
    features = analyze_urgency_features(tokens)
    active_tokens = features["active_tokens"]
    critical_keywords = features["critical_keywords"]
    medium_keywords = features["medium_keywords"]
    high_academic_impact = bool(features["high_impact"])
    operational_block = bool(features["operational_block"])

    # Safety evidence is authoritative and cannot be lowered by positive
    # sentiment or by the MNB prediction.
    if critical_keywords:
        logger.info("Urgensi tinggi karena indikator keselamatan: %s", critical_keywords)
        return "Tinggi"
    if high_academic_impact:
        logger.info("Urgensi tinggi karena dampak akademik kritis")
        return "Tinggi"

    has_operational_issue = bool(medium_keywords) or operational_block
    negative_support = sentiment_score < 0
    rule_urgency = "Sedang" if has_operational_issue or negative_support else "Rendah"

    model_urgency = predict_urgency_from_model(active_tokens)
    if model_urgency == "Tinggi":
        # MNB is not allowed to create a high class.  It can only support a
        # medium result when the text itself shows a real issue.
        return "Sedang" if has_operational_issue or negative_support else "Rendah"
    if model_urgency == "Sedang" and has_operational_issue:
        return "Sedang"
    if model_urgency == "Rendah" and rule_urgency == "Sedang":
        return "Sedang"
    if model_urgency in URGENCY_LEVELS and rule_urgency == "Rendah":
        return model_urgency if model_urgency == "Rendah" else "Rendah"
    return rule_urgency
