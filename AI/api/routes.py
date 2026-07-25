"""Route API untuk modul prediksi AI."""

from __future__ import annotations

import logging
from typing import Any

from fastapi import APIRouter, HTTPException, status
from pydantic import BaseModel, Field, field_validator

from services.prediction import predict_complaint

logger = logging.getLogger(__name__)

router = APIRouter()


class PredictRequest(BaseModel):
    """Schema request endpoint prediksi."""

    deskripsi: str = Field(..., min_length=1)

    @field_validator("deskripsi")
    @classmethod
    def deskripsi_not_blank(cls, value: str) -> str:
        """Reject descriptions that contain only whitespace."""
        if not value.strip():
            raise ValueError("Deskripsi wajib diisi")
        return value


class PredictResponse(BaseModel):
    """Schema response endpoint prediksi."""

    cleaned_text: str
    tokens: list[str]
    score: int
    sentiment: str
    urgency: str


@router.post("/predict", response_model=PredictResponse, status_code=status.HTTP_200_OK)
def predict(payload: PredictRequest) -> dict[str, Any]:
    """Prediksi sentimen dan urgensi dari deskripsi pengaduan."""
    try:
        return predict_complaint(payload.deskripsi)
    except ValueError as exc:
        logger.warning("Request prediksi tidak valid: %s", exc)
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=str(exc)) from exc
    except Exception as exc:
        logger.exception("Terjadi kesalahan saat prediksi")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Terjadi kesalahan pada server",
        ) from exc
