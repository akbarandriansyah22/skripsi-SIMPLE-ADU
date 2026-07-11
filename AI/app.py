"""Aplikasi FastAPI untuk modul AI pengaduan mahasiswa."""

from __future__ import annotations

import logging

from fastapi import FastAPI

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
    """Health check sederhana untuk memastikan API aktif."""
    return {"status": "ok"}
