"""Unit test sederhana untuk endpoint /predict."""

from __future__ import annotations

import pytest
from fastapi.exceptions import RequestValidationError
from pydantic import ValidationError

from api.routes import PredictRequest, predict


def _raise_request_validation_error(error: ValidationError) -> None:
    """Samakan error validasi Pydantic dengan respons validasi FastAPI."""
    raise RequestValidationError(error.errors()) from error


def test_predict_endpoint_success() -> None:
    """Endpoint /predict mengembalikan struktur response yang benar."""
    payload = PredictRequest(deskripsi="Fasilitas kampus rusak dan harus segera diperbaiki.")
    data = predict(payload)

    assert set(data.keys()) == {"cleaned_text", "tokens", "score", "sentiment", "urgency"}
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
    assert data["score"] == 0
    assert data["sentiment"] == "Netral"
    assert data["urgency"] == "Sedang"


def test_predict_endpoint_without_urgency_keyword() -> None:
    payload = PredictRequest(deskripsi="Informasi jadwal kuliah semester ini belum jelas.")
    data = predict(payload)

    assert data["urgency"] in {"Rendah", "Sedang", "Tinggi"}
    assert isinstance(data["score"], int)
