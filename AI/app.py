"""Aplikasi FastAPI untuk modul AI pengaduan mahasiswa."""

from __future__ import annotations

import logging

from fastapi import FastAPI, status
from fastapi.responses import JSONResponse

from api.routes import router

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(name)s - %(message)s",
)

app = FastAPI(
    title="Simpel ADU AI",
    description="API rule-based untuk analisis sentimen dan urgensi pengaduan mahasiswa.",
    version="1.0.0",
)

app.include_router(router)


@app.get("/health")
def health_check() -> dict[str, str]:
    """Readiness check verifies the model and both sentiment lexicons."""
    try:
        from services.sentiment import load_sentiment_lexicons
        from services.urgency_model import load_urgency_model

        load_sentiment_lexicons()
        if load_urgency_model() is None:
            raise RuntimeError("model urgensi tidak tersedia atau rusak")
        return {"status": "ok"}
    except Exception:
        logging.getLogger(__name__).exception("AI belum siap")
        return JSONResponse(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            content={"status": "not_ready"},
        )
