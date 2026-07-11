"""Utility untuk membaca dan memvalidasi dataset Excel."""

from __future__ import annotations

import logging
from pathlib import Path

import pandas as pd

logger = logging.getLogger(__name__)

DEFAULT_DATASET_PATH = Path(__file__).resolve().parents[1] / "dataset" / "dataset.xlsx"


def read_dataset(file_path: str | Path = DEFAULT_DATASET_PATH) -> pd.DataFrame:
    """Baca dataset Excel dan tampilkan ringkasan kualitas data.

    Args:
        file_path: Lokasi file Excel yang akan dibaca.

    Returns:
        DataFrame berisi data dari file Excel.

    Raises:
        FileNotFoundError: Jika file tidak ditemukan.
        ValueError: Jika file kosong atau gagal dibaca.
    """
    path = Path(file_path)
    logger.info("Membaca dataset dari %s", path)

    if not path.exists():
        logger.error("Dataset tidak ditemukan: %s", path)
        raise FileNotFoundError(f"Dataset tidak ditemukan: {path}")

    if path.suffix.lower() not in {".xlsx", ".xls"}:
        logger.error("Format file dataset tidak valid: %s", path.suffix)
        raise ValueError("Dataset harus berupa file Excel (.xlsx atau .xls)")

    try:
        dataframe = pd.read_excel(path)
    except Exception as exc:
        logger.exception("Gagal membaca dataset Excel")
        raise ValueError(f"Gagal membaca dataset Excel: {exc}") from exc

    if dataframe.empty:
        logger.warning("Dataset kosong: %s", path)
        raise ValueError("Dataset kosong")

    print(f"Jumlah data: {len(dataframe)}")
    print(f"Nama kolom: {list(dataframe.columns)}")
    print("Missing value:")
    print(dataframe.isna().sum())
    print("5 data pertama:")
    print(dataframe.head())

    return dataframe


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    read_dataset()
