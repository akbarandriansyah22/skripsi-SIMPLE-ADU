"""Review the active questionnaire dataset without modifying the source file."""

from __future__ import annotations

import json
from collections import Counter
from pathlib import Path

import pandas as pd

from services.preprocessing import preprocess_text
from services.urgency import analyze_urgency_features
from utils.read_excel import DEFAULT_DATASET_PATH, read_dataset

REVIEWED_DATASET_PATH = Path(__file__).resolve().parent / "data" / "dataset_urgency_reviewed.xlsx"
AUDIT_REPORT_PATH = Path(__file__).resolve().parent / "reports" / "urgency_dataset_audit.md"

TEXT_COLUMN = "Tuliskan pengaduan, keluhan, atau saran yang pernah Anda alami selama menjadi mahasiswa FT UMJ.\n"
CATEGORY_COLUMN = "Pilih kategori yang paling sesuai dengan pengaduan Anda"
CONDITION_COLUMN = "Bagaimana kondisi yang Anda alami terkait pengaduan tersebut?"

CONDITION_TO_URGENCY = {
    "Biasa saja": "Rendah",
    "Tidak terlalu bermasalah": "Rendah",
    "Cukup mengganggu": "Sedang",
    "Sangat merugikan": "Tinggi",
}

# Contextual decisions are deliberately explicit so every reviewed row is
# auditable.  The source condition is retained only as a supporting feature.
REVIEW_DECISIONS: dict[int, tuple[str, str, str, bool]] = {
    1: ("Sedang", "high", "Gangguan kebersihan nyata pada fasilitas; tidak ada risiko keselamatan langsung.", False),
    2: ("Rendah", "medium", "Kekurangan media pengaduan merupakan kebutuhan perbaikan sistem, bukan gangguan layanan yang sedang berlangsung.", False),
    3: ("Sedang", "medium", "Jam layanan fasilitas tidak konsisten dan pengaduan tidak ditindaklanjuti; dampak nyata tetapi tanpa kondisi kritis.", False),
    4: ("Sedang", "medium", "Akses komunikasi dengan dosen terhambat dan berkaitan dengan perbaikan nilai, tetapi belum ada bukti kehilangan hak akademik.", False),
    5: ("Sedang", "medium", "Perkuliahan tidak berjalan optimal; aktivitas masih dapat dilanjutkan dan tidak ada bahaya langsung.", False),
    6: ("Rendah", "high", "Permintaan penambahan alat gym; tidak menyatakan gangguan layanan atau kondisi kritis.", False),
    7: ("Rendah", "high", "Respon menyatakan tidak ada masalah.", False),
    8: ("Sedang", "high", "AC mati membuat ruang panas dan mengganggu perkuliahan, tetapi tidak membahayakan secara langsung.", False),
    9: ("Sedang", "medium", "Kekurangan dosen, infrastruktur, dan kerja sama menghambat kualitas layanan tanpa indikator kritis.", False),
    10: ("Sedang", "high", "Tidak dapat login adalah gangguan akses layanan administrasi yang masih dapat ditindaklanjuti; tidak ada bukti dampak kritis.", False),
    11: ("Sedang", "high", "Beberapa gangguan akademik dan fasilitas nyata disebutkan; dampak luas atau batas waktu kritis tidak dibuktikan.", False),
    12: ("Rendah", "low", "Teks terlalu singkat dan tidak menjelaskan gangguan administrasi yang spesifik; digunakan label konservatif.", True),
    13: ("Rendah", "low", "Teks singkat tentang kekurangan fasilitas tanpa dampak aktivitas yang jelas; digunakan label konservatif.", True),
    14: ("Rendah", "low", "Kekurangan informasi akademik disebutkan tanpa contoh dampak; belum cukup untuk urgensi lebih tinggi.", True),
    15: ("Rendah", "medium", "Respon belum ada tidak menjelaskan masalah aktif yang perlu penanganan segera.", False),
    16: ("Sedang", "high", "AC bermasalah di kelas merupakan gangguan fasilitas nyata, bukan bahaya langsung.", False),
    17: ("Sedang", "medium", "Batas pendaftaran beasiswa merupakan kendala kebijakan yang perlu ditinjau, tetapi tidak membuktikan keadaan darurat.", False),
    18: ("Sedang", "high", "Informasi jadwal susulan yang mendadak menghambat persiapan, tanpa bukti batas waktu kehilangan hak.", False),
    19: ("Sedang", "medium", "Kepadatan kelas dan informasi mendadak mengganggu kegiatan, tetapi tidak menunjukkan kondisi kritis.", False),
    20: ("Rendah", "high", "Respon menyatakan tidak ada masalah.", False),
    21: ("Sedang", "high", "Sistem terkadang down; gangguan operasional berulang tetapi layanan tidak dinyatakan berhenti total.", False),
    22: ("Sedang", "medium", "Akses parkir yang jauh menghambat kenyamanan/akses, tanpa risiko keselamatan atau batas waktu kritis.", False),
    23: ("Sedang", "low", "Nilai konversi dan e-learning disebut sebagai kendala, tetapi detail dampaknya tidak dijelaskan.", True),
    24: ("Rendah", "high", "Teks menyatakan fasilitas masih banyak tersedia dan tidak menunjukkan gangguan.", False),
    25: ("Sedang", "high", "Keterbatasan laboratorium adalah kendala fasilitas yang dapat menghambat aktivitas, bukan kondisi kritis.", False),
    26: ("Rendah", "low", "Teks hanya menyebut sistem kurikulum tanpa masalah atau dampak yang dapat dinilai.", True),
    27: ("Rendah", "low", "Teks hanya menyebut masalah toilet tanpa menjelaskan kondisi atau dampaknya.", True),
    28: ("Rendah", "low", "Teks hanya menyebut lapangan futsal; belum ada keluhan atau dampak yang jelas.", True),
    29: ("Rendah", "high", "Permintaan pembaruan laboratorium adalah saran pengembangan, bukan gangguan kritis.", False),
    30: ("Rendah", "low", "Teks hanya menyebut sistem akademik tanpa konteks gangguan.", True),
    31: ("Rendah", "high", "Permintaan penambahan laboratorium merupakan saran kapasitas, bukan kondisi darurat.", False),
    32: ("Rendah", "high", "Respon menyatakan tidak ada masalah.", False),
    33: ("Rendah", "low", "Teks hanya menyebut masalah organisasi tanpa dampak layanan yang spesifik.", True),
    34: ("Sedang", "high", "Toilet kotor adalah gangguan fasilitas nyata tanpa indikator keselamatan/kesehatan kritis.", False),
    35: ("Rendah", "high", "Keterbatasan buku merupakan kebutuhan fasilitas dengan dampak ringan; aktivitas utama tetap dapat berlangsung.", False),
    36: ("Sedang", "high", "Banyak kursi rusak mengganggu fasilitas kelas, tetapi tidak ada bukti bahaya langsung.", False),
    37: ("Sedang", "high", "Lapangan tidak dapat dipakai merupakan gangguan fasilitas tanpa cakupan atau batas waktu kritis.", False),
    38: ("Rendah", "low", "Teks hanya menyebut AC di ruang kelas tanpa menjelaskan kerusakan atau dampak.", True),
    39: ("Sedang", "high", "AC berisik mengganggu kenyamanan kelas, tetapi layanan masih berjalan.", False),
    40: ("Sedang", "high", "Dosen kadang terlambat; gangguan operasional berulang tanpa dampak akademik kritis.", False),
    41: ("Sedang", "high", "Nilai terlambat keluar menghambat administrasi mahasiswa, tetapi tidak menyebut batas akhir atau kehilangan hak.", False),
    42: ("Rendah", "high", "Kegiatan organisasi yang jarang adalah keluhan ringan tanpa hambatan utama.", False),
    43: ("Rendah", "high", "Permintaan media penghubung adalah saran pengembangan sistem.", False),
    44: ("Rendah", "high", "Permintaan kemudahan keuangan untuk KRS adalah saran; tidak ada gangguan aktif yang dijelaskan.", False),
    45: ("Rendah", "high", "Permintaan kotak aspirasi/media pengaduan adalah saran fasilitas komunikasi.", False),
    46: ("Rendah", "low", "Teks hanya menyebut mata kuliah kesenian tanpa masalah atau dampak yang jelas.", True),
    47: ("Sedang", "high", "Dosen kadang terlambat memberi informasi dan mengganggu proses, tetapi tanpa batas waktu kritis.", False),
    48: ("Rendah", "high", "Permintaan penambahan komputer laboratorium adalah saran kapasitas.", False),
    49: ("Sedang", "high", "Toilet dan AC kadang tidak berfungsi merupakan gangguan fasilitas berulang, tetapi masih tanpa kondisi kritis.", False),
    50: ("Sedang", "high", "Keterlambatan dosen mengganggu kegiatan, tetapi tidak membuktikan dampak kritis.", False),
    51: ("Sedang", "high", "Proposal acara lama ditandatangani; proses administratif terhambat tanpa batas waktu kritis yang disebutkan.", False),
    52: ("Sedang", "high", "Keluhan penilaian dosen perlu ditindaklanjuti, tetapi tidak ada bukti kehilangan hak akademik.", False),
    53: ("Sedang", "high", "Pengisian KRS terlambat diproses dosen; ini gangguan akademik, namun tidak ada batas akhir/kegagalan nyata.", False),
    54: ("Rendah", "high", "Permintaan penambahan komputer laboratorium adalah saran kapasitas, bukan gangguan kritis.", False),
}


def _as_text(value: object) -> str:
    return "" if pd.isna(value) else str(value).strip()


def _initial_label(condition: str) -> str:
    return CONDITION_TO_URGENCY.get(condition, "Tidak valid")


def build_reviewed_dataset(dataframe: pd.DataFrame) -> pd.DataFrame:
    rows: list[dict[str, object]] = []
    for position, (_, row) in enumerate(dataframe.iterrows(), start=1):
        text = _as_text(row.get(TEXT_COLUMN))
        category = _as_text(row.get(CATEGORY_COLUMN))
        condition = _as_text(row.get(CONDITION_COLUMN))
        initial = _initial_label(condition)

        if position not in REVIEW_DECISIONS:
            raise ValueError(f"Keputusan review belum tersedia untuk nomor {position}")
        reviewed, confidence, reason, manual = REVIEW_DECISIONS[position]

        tokens = preprocess_text(text)["tokens"] if text else []
        features = analyze_urgency_features(tokens)
        changed_reason = reason if initial != reviewed else "Label dipertahankan setelah pemeriksaan konteks teks. " + reason
        rows.append(
            {
                "nomor_data": position,
                "teks_pengaduan": text,
                "kategori_asli": category,
                "kondisi_asli": condition,
                "label_awal": initial,
                "label_reviewed": reviewed if text and initial != "Tidak valid" else "Rendah",
                "confidence": "low" if manual or not text else confidence,
                "alasan_review": "Data kosong/tidak valid." if not text else changed_reason,
                "critical_keywords": ", ".join(features["critical_keywords"]),
                "high_impact_indicators": ", ".join(
                    name for name, present in features["high_impact_indicators"].items() if present
                ),
                "medium_keywords": ", ".join(features["medium_keywords"]),
                "mitigating_phrases": ", ".join(features["mitigating_phrases"]),
                "negation_detected": bool(features["negation_detected"]),
                "manual_review_required": bool(manual or not text),
            }
        )

    return pd.DataFrame(rows)


def _distribution(values: pd.Series, labels: tuple[str, ...] = ("Rendah", "Sedang", "Tinggi")) -> str:
    counts = Counter(str(value) for value in values)
    return ", ".join(f"{label}: {counts.get(label, 0)}" for label in labels)


def build_audit_report(source: pd.DataFrame, reviewed: pd.DataFrame) -> str:
    changed = reviewed[reviewed["label_awal"] != reviewed["label_reviewed"]]
    ambiguous = reviewed[reviewed["manual_review_required"]]
    valid = reviewed[reviewed["label_awal"] != "Tidak valid"]

    lines = [
        "# Audit Dataset Urgensi SIMPEL-ADU",
        "",
        f"- Dataset asli: `{DEFAULT_DATASET_PATH.relative_to(DEFAULT_DATASET_PATH.parents[1])}`",
        f"- Dataset asli tetap utuh; hasil review: `{REVIEWED_DATASET_PATH.relative_to(REVIEWED_DATASET_PATH.parents[1])}`",
        f"- Jumlah data asli: {len(source)}",
        f"- Baris valid: {len(valid)}",
        f"- Baris kosong/tidak valid: {len(reviewed) - len(valid)}",
        "",
        "## Distribusi",
        "",
        f"- Kondisi responden: {_distribution(source[CONDITION_COLUMN].fillna(''), tuple(CONDITION_TO_URGENCY))}",
        f"- Label awal (pemetaan lama): {_distribution(reviewed['label_awal'])}",
        f"- Label reviewed: {_distribution(reviewed['label_reviewed'])}",
        f"- Jumlah label berubah: {len(changed)}",
        f"- Data manual review: {len(ambiguous)}",
        f"- Data Tinggi yang benar-benar valid: {sum(reviewed['label_reviewed'] == 'Tinggi')}",
        "- Data referensi curated untuk pelatihan: 0 (tidak digunakan).",
        "",
        "## Pemisahan sentimen dan urgensi",
        "",
        "`negative.tsv` dan `positive.tsv` tetap menjadi lexicon polaritas sentimen. Skor sentimen tidak dipakai sebagai pemetaan langsung ke urgensi. Sentimen negatif hanya dapat mendukung kenaikan Rendah menjadi Sedang; Tinggi memerlukan bukti keselamatan atau dampak akademik kritis.",
        "",
        "## Data yang labelnya berubah",
        "",
        "| Nomor | Teks | Kondisi asli | Awal | Reviewed | Alasan |",
        "|---:|---|---|---|---|---|",
    ]
    for _, row in changed.iterrows():
        text = str(row["teks_pengaduan"]).replace("|", "\\|").replace("\n", " ")
        reason = str(row["alasan_review"]).replace("|", "\\|")
        lines.append(
            f"| {row['nomor_data']} | {text} | {row['kondisi_asli']} | {row['label_awal']} | {row['label_reviewed']} | {reason} |"
        )

    lines.extend(["", "## Data ambigu/manual review", ""])
    for _, row in ambiguous.iterrows():
        lines.append(f"- {row['nomor_data']}: {row['teks_pengaduan']} — {row['alasan_review']}")

    lines.extend(
        [
            "",
            "## Catatan evaluasi",
            "",
            "Tidak ada respons kuesioner yang memenuhi bukti Tinggi secara kuat. Karena itu tidak ada label Tinggi yang dipaksakan dan tidak ada data kuesioner yang diubah menjadi Tinggi untuk menyeimbangkan kelas. Deteksi Tinggi pada input baru tetap dijaga oleh guardrail keselamatan dan pola dampak akademik kritis.",
            "",
            "Dataset ini tidak digunakan dalam proses AI aktif; urgensi final ditentukan dengan rule-based.",
            "",
        ]
    )
    return "\n".join(lines)


def review_dataset() -> tuple[pd.DataFrame, str]:
    source = read_dataset(DEFAULT_DATASET_PATH)
    reviewed = build_reviewed_dataset(source)
    REVIEWED_DATASET_PATH.parent.mkdir(parents=True, exist_ok=True)
    AUDIT_REPORT_PATH.parent.mkdir(parents=True, exist_ok=True)
    reviewed.to_excel(REVIEWED_DATASET_PATH, index=False)
    report = build_audit_report(source, reviewed)
    AUDIT_REPORT_PATH.write_text(report, encoding="utf-8")
    return reviewed, report


if __name__ == "__main__":
    reviewed, _ = review_dataset()
    print(json.dumps({
        "rows": len(reviewed),
        "label_counts": dict(Counter(reviewed["label_reviewed"])),
        "changed": int((reviewed["label_awal"] != reviewed["label_reviewed"]).sum()),
        "manual_review": int(reviewed["manual_review_required"].sum()),
        "high_valid": int((reviewed["label_reviewed"] == "Tinggi").sum()),
    }, ensure_ascii=False, indent=2))
