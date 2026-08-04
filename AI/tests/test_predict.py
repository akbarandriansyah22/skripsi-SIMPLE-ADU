"""Unit test sederhana untuk endpoint /predict."""

from __future__ import annotations

import pytest
from fastapi.exceptions import RequestValidationError
from pydantic import ValidationError

from api.routes import PredictRequest, predict
from app import health_check
from services.prediction import predict_complaint
from services import sentiment


def _raise_request_validation_error(error: ValidationError) -> None:
    """Samakan error validasi Pydantic dengan respons validasi FastAPI."""
    raise RequestValidationError(error.errors()) from error


def test_predict_http_endpoint() -> None:
    """Route handler mempertahankan kontrak response endpoint /predict."""
    body = predict(PredictRequest(deskripsi="Ada kebakaran di laboratorium."))
    assert {"cleaned_text", "tokens", "score", "sentiment", "urgency"}.issubset(body)
    assert {"original_text", "positive_score", "negative_score", "sentiment_score", "sentiment_label", "matched_words", "urgency_score", "urgency_label", "urgency_reason"}.issubset(body)
    assert body["urgency"] == "Tinggi"


def test_health_http_endpoint() -> None:
    assert health_check() == {"status": "ok"}


def test_predict_endpoint_success() -> None:
    """Endpoint /predict mengembalikan struktur response yang benar."""
    payload = PredictRequest(deskripsi="Fasilitas kampus rusak dan harus segera diperbaiki.")
    data = predict(payload)

    assert {"cleaned_text", "tokens", "score", "sentiment", "urgency"}.issubset(data)
    assert isinstance(data["cleaned_text"], str)
    assert isinstance(data["tokens"], list)
    assert isinstance(data["score"], int)
    assert data["sentiment"] in {"Positif", "Negatif", "Netral"}
    assert data["urgency"] in {"Rendah", "Sedang", "Tinggi"}


def test_predict_endpoint_empty_description() -> None:
    """Endpoint /predict menolak deskripsi kosong."""
    with pytest.raises(RequestValidationError):
        try:
            PredictRequest(deskripsi="")
        except ValidationError as exc:
            _raise_request_validation_error(exc)


def test_predict_endpoint_blank_description() -> None:
    """Endpoint /predict menolak deskripsi berisi spasi."""
    with pytest.raises(RequestValidationError):
        try:
            PredictRequest(deskripsi="   ")
        except ValidationError as exc:
            _raise_request_validation_error(exc)


def test_predict_endpoint_fire_complaint_is_high_urgency() -> None:
    """Pengaduan kebakaran harus masuk urgensi tinggi."""
    payload = PredictRequest(
        deskripsi="Terjadi kebakaran di laboratorium komputer dan asap sudah memenuhi ruangan."
    )
    data = predict(payload)

    assert data["urgency"] == "Tinggi"


def test_predict_endpoint_negative_sentiment() -> None:
    payload = PredictRequest(deskripsi="Pelayanan buruk dan fasilitas rusak parah.")
    data = predict(payload)

    assert data["sentiment"] in {"Negatif", "Netral"}


def test_predict_endpoint_neutral_text() -> None:
    payload = PredictRequest(deskripsi="Saya ingin menanyakan jadwal penggunaan laboratorium.")
    data = predict(payload)

    assert data["sentiment"] in {"Positif", "Negatif", "Netral"}


def test_predict_endpoint_medium_urgency_keyword() -> None:
    payload = PredictRequest(deskripsi="AC ruang kelas rusak dan perlu diperbaiki.")
    data = predict(payload)

    assert data["urgency"] in {"Sedang", "Tinggi"}


def test_predict_endpoint_wifi_lab_complaint_is_medium_urgency() -> None:
    payload = PredictRequest(
        deskripsi="Wifi laboratorium lambat dan mengganggu kegiatan praktikum."
    )
    data = predict(payload)

    assert data["tokens"] == ["wifi", "laboratorium", "lambat", "ganggu", "giat", "praktikum"]
    # The project lexicon intentionally scores both "lambat" and "ganggu"
    # as negative terms; urgency remains medium because both are medium-urgency
    # keywords.
    assert data["score"] == -6
    assert data["sentiment"] == "Negatif"
    assert data["urgency"] == "Sedang"


def test_predict_endpoint_without_urgency_keyword() -> None:
    payload = PredictRequest(deskripsi="Informasi jadwal kuliah semester ini belum jelas.")
    data = predict(payload)

    assert data["urgency"] in {"Rendah", "Sedang", "Tinggi"}
    assert isinstance(data["score"], int)


@pytest.mark.parametrize(
    ("description", "must_be_high"),
    [
        ("Pelayanan buruk.", False),
        ("Pelayanan tidak buruk.", False),
        ("Fasilitas aman.", False),
        ("Fasilitas tidak aman.", False),
        ("Masalah belum diperbaiki.", False),
        ("Tidak ada kebakaran.", False),
        ("Ada kebakaran.", True),
        ("AC tidak dingin.", False),
        ("Jaringan tidak lambat.", False),
        ("Jaringan sangat lambat.", False),
    ],
)
def test_negation_and_urgency_cases(description: str, must_be_high: bool) -> None:
    """Negation is retained and must not create a false critical escalation."""
    data = predict(PredictRequest(deskripsi=description))

    if "tidak" in description.lower() or "belum" in description.lower():
        assert any(token in data["tokens"] for token in ("tidak", "belum"))
    assert (data["urgency"] == "Tinggi") is must_be_high


def test_negation_reverses_known_lexicon_polarity(monkeypatch: pytest.MonkeyPatch) -> None:
    """The polarity inversion is bounded to the next sentiment-bearing term."""
    monkeypatch.setattr(
        sentiment,
        "load_sentiment_lexicons",
        lambda: ({"aman": 3}, {"buruk": -3}),
    )

    assert sentiment.analyze_sentiment(["buruk"])["score"] == -3
    assert sentiment.analyze_sentiment(["tidak", "buruk"])["score"] == 3
    assert sentiment.analyze_sentiment(["aman"])["score"] == 3
    assert sentiment.analyze_sentiment(["tidak", "aman"])["score"] == -3


@pytest.mark.parametrize(
    ("description", "expected_urgency"),
    [
        (
            "Koneksi Wi-Fi di laboratorium sering lambat dan mengganggu kegiatan praktikum. "
            "Mahasiswa mengalami kesulitan mengakses materi pembelajaran dan mengunggah tugas, "
            "tetapi jaringan masih dapat digunakan.",
            "Sedang",
        ),
        ("Terjadi korsleting listrik di laboratorium komputer dan muncul asap dari stop kontak. Kondisi ini membahayakan mahasiswa.", "Tinggi"),
        ("Saya ingin menanyakan jadwal penggunaan laboratorium komputer.", "Rendah"),
        ("Gym tolong ditambahkan lagi alatnya.", "Rendah"),
        ("AC mati sehingga ruang kelas panas dan mengganggu kegiatan perkuliahan.", "Sedang"),
        ("Sistem terkadang down, tetapi masih dapat diakses kembali setelah beberapa saat.", "Sedang"),
        ("Lab informatika sebaiknya diperbanyak komputernya.", "Rendah"),
        ("Nilai mata kuliah terlambat keluar dan menghambat proses administrasi mahasiswa.", "Sedang"),
        ("Tidak terjadi korsleting dan tidak ada asap di laboratorium.", "Rendah"),
        (
            "Sistem KRS tidak dapat diakses sama sekali oleh seluruh mahasiswa dan hari ini batas akhir pengisian. "
            "Tidak ada alternatif sehingga mahasiswa terancam gagal mengambil mata kuliah.",
            "Tinggi",
        ),
        (
            "Pelayanan sangat lambat, sulit, buruk, dan mengganggu, tetapi proses masih dapat dilakukan.",
            "Sedang",
        ),
    ],
)
def test_required_urgency_examples(description: str, expected_urgency: str) -> None:
    data = predict_complaint(description)
    assert data["urgency"] == expected_urgency


def test_required_wifi_example_has_negative_sentiment() -> None:
    data = predict_complaint(
        "Koneksi Wi-Fi di laboratorium sering lambat dan mengganggu kegiatan praktikum. "
        "Mahasiswa mengalami kesulitan mengakses materi pembelajaran dan mengunggah tugas, "
        "tetapi jaringan masih dapat digunakan."
    )
    assert data["sentiment"] == "Negatif"


def test_very_negative_noncritical_text_cannot_be_high() -> None:
    data = predict_complaint(
        "Pelayanan sangat lambat, sulit, buruk, dan mengganggu, tetapi proses masih dapat dilakukan."
    )
    assert data["sentiment"] != "Positif"
    assert data["urgency"] == "Sedang"
