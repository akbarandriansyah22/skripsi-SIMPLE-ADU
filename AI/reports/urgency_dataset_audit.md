# Audit Dataset Urgensi SIMPEL-ADU

- Dataset asli: `dataset/dataset.xlsx`
- Dataset asli tetap utuh; hasil review: `data/dataset_urgency_reviewed.xlsx`
- Jumlah data asli: 54
- Baris valid: 54
- Baris kosong/tidak valid: 0

## Distribusi

- Kondisi responden: Biasa saja: 7, Tidak terlalu bermasalah: 4, Cukup mengganggu: 27, Sangat merugikan: 16
- Label awal (pemetaan lama): Rendah: 11, Sedang: 27, Tinggi: 16
- Label reviewed: Rendah: 26, Sedang: 28, Tinggi: 0
- Jumlah label berubah: 28
- Data manual review: 11
- Data Tinggi yang benar-benar valid: 0
- Data referensi curated untuk pelatihan: 0 (tidak digunakan).

## Pemisahan sentimen dan urgensi

`negative.tsv` dan `positive.tsv` tetap menjadi lexicon polaritas sentimen. Skor sentimen tidak dipakai sebagai pemetaan langsung ke urgensi. Sentimen negatif hanya dapat mendukung kenaikan Rendah menjadi Sedang; Tinggi memerlukan bukti keselamatan atau dampak akademik kritis.

## Data yang labelnya berubah

| Nomor | Teks | Kondisi asli | Awal | Reviewed | Alasan |
|---:|---|---|---|---|---|
| 3 | masalah fasilitas tuh ada fasilitas kmpus yang kenapa ketika buka dan tutup yang saya alami dulu ya beberapa bulan yg lalu itu tutup dan buka sesuai sibuk nya tendik mksdnya ini kan fasilitas ya kenapa malah mengikuti jadwal mereka, saya ga tau ini masih terjadi pada saat ini tpi ini pernah terjadi bulan lalu pas puasa dan pengaduan nya tidak di respon baik juga mksdnya ga ada kabar kelanjutan nya gmn | Sangat merugikan | Tinggi | Sedang | Jam layanan fasilitas tidak konsisten dan pengaduan tidak ditindaklanjuti; dampak nyata tetapi tanpa kondisi kritis. |
| 4 | ada dosen yang tidak mau bertemu dengan mahasiswa, padahal mahasiswa sendiri mau meminta ke mudahan atau kesempatan untuk memperbaiki nilai | Sangat merugikan | Tinggi | Sedang | Akses komunikasi dengan dosen terhambat dan berkaitan dengan perbaikan nilai, tetapi belum ada bukti kehilangan hak akademik. |
| 5 | Pembelajaran mata kuliah yang offline tapi ga Dateng giliran Dateng cuma bentar dan itu cuma jelasin materi ppt | Sangat merugikan | Tinggi | Sedang | Perkuliahan tidak berjalan optimal; aktivitas masih dapat dilanjutkan dan tidak ada bahaya langsung. |
| 6 | Gym tolong di tambahkan lagi alatnya | Cukup mengganggu | Sedang | Rendah | Permintaan penambahan alat gym; tidak menyatakan gangguan layanan atau kondisi kritis. |
| 12 | Administrasi mostly. | Cukup mengganggu | Sedang | Rendah | Teks terlalu singkat dan tidak menjelaskan gangguan administrasi yang spesifik; digunakan label konservatif. |
| 13 | fasilitas kelas ada yang kurang | Cukup mengganggu | Sedang | Rendah | Teks singkat tentang kekurangan fasilitas tanpa dampak aktivitas yang jelas; digunakan label konservatif. |
| 15 | Belum ada | Cukup mengganggu | Sedang | Rendah | Respon belum ada tidak menjelaskan masalah aktif yang perlu penanganan segera. |
| 16 | kendala AC yang bermasalah di kelas ibnu firnas | Sangat merugikan | Tinggi | Sedang | AC bermasalah di kelas merupakan gangguan fasilitas nyata, bukan bahaya langsung. |
| 17 | Mengenai beasiswa KIP yang hanya bisa didaftarkan hanya pada saat mendaftar mahasiswa baru | Sangat merugikan | Tinggi | Sedang | Batas pendaftaran beasiswa merupakan kendala kebijakan yang perlu ditinjau, tetapi tidak membuktikan keadaan darurat. |
| 19 | Pada selama kuliah, karena mahasiswa terlalu banyak pada saat melihat persentase yang dijelaskan dari dosen bangku kelas tidak kelihatan, dan penjadwalan mata kuliah terkadang suka diberikan informasi mendadak. | Sangat merugikan | Tinggi | Sedang | Kepadatan kelas dan informasi mendadak mengganggu kegiatan, tetapi tidak menunjukkan kondisi kritis. |
| 24 | Fasilitas nya masih banyak tersedia | Cukup mengganggu | Sedang | Rendah | Teks menyatakan fasilitas masih banyak tersedia dan tidak menunjukkan gangguan. |
| 27 | Masalah toilet | Cukup mengganggu | Sedang | Rendah | Teks hanya menyebut masalah toilet tanpa menjelaskan kondisi atau dampaknya. |
| 29 | Lab nya diperbarui | Cukup mengganggu | Sedang | Rendah | Permintaan pembaruan laboratorium adalah saran pengembangan, bukan gangguan kritis. |
| 30 | Sistem akademik | Sangat merugikan | Tinggi | Rendah | Teks hanya menyebut sistem akademik tanpa konteks gangguan. |
| 31 | Labnya di perbanyak | Sangat merugikan | Tinggi | Rendah | Permintaan penambahan laboratorium merupakan saran kapasitas, bukan kondisi darurat. |
| 32 | Tidak ada | Sangat merugikan | Tinggi | Rendah | Respon menyatakan tidak ada masalah. |
| 36 | Kursinya masih banyak yg rusak | Sangat merugikan | Tinggi | Sedang | Banyak kursi rusak mengganggu fasilitas kelas, tetapi tidak ada bukti bahaya langsung. |
| 37 | Lapangan futsal belum bisa dipakai | Biasa saja | Rendah | Sedang | Lapangan tidak dapat dipakai merupakan gangguan fasilitas tanpa cakupan atau batas waktu kritis. |
| 38 | Ac diruang kelas | Sangat merugikan | Tinggi | Rendah | Teks hanya menyebut AC di ruang kelas tanpa menjelaskan kerusakan atau dampak. |
| 41 | Nilai terlambat keluar | Sangat merugikan | Tinggi | Sedang | Nilai terlambat keluar menghambat administrasi mahasiswa, tetapi tidak menyebut batas akhir atau kehilangan hak. |
| 44 | Adanya kemudahan keuangan buat isi krs | Cukup mengganggu | Sedang | Rendah | Permintaan kemudahan keuangan untuk KRS adalah saran; tidak ada gangguan aktif yang dijelaskan. |
| 45 | Adanya kotak aspirasi di fakultas ata sistem lainya | Cukup mengganggu | Sedang | Rendah | Permintaan kotak aspirasi/media pengaduan adalah saran fasilitas komunikasi. |
| 46 | Terkait mk kesenian | Cukup mengganggu | Sedang | Rendah | Teks hanya menyebut mata kuliah kesenian tanpa masalah atau dampak yang jelas. |
| 47 | Dosen kadang suka telat kasih info | Sangat merugikan | Tinggi | Sedang | Dosen kadang terlambat memberi informasi dan mengganggu proses, tetapi tanpa batas waktu kritis. |
| 48 | Lab mesin masih harus di perbanyak | Cukup mengganggu | Sedang | Rendah | Permintaan penambahan komputer laboratorium adalah saran kapasitas. |
| 51 | Perihal proposal acara yang suka lama di TTD oleh wadek | Sangat merugikan | Tinggi | Sedang | Proposal acara lama ditandatangani; proses administratif terhambat tanpa batas waktu kritis yang disebutkan. |
| 53 | Pengisian KRS kadang suka telat dari dospen | Sangat merugikan | Tinggi | Sedang | Pengisian KRS terlambat diproses dosen; ini gangguan akademik, namun tidak ada batas akhir/kegagalan nyata. |
| 54 | Lab informatika harus diperbanyak komputer nya | Sangat merugikan | Tinggi | Rendah | Permintaan penambahan komputer laboratorium adalah saran kapasitas, bukan gangguan kritis. |

## Data ambigu/manual review

- 12: Administrasi mostly. — Teks terlalu singkat dan tidak menjelaskan gangguan administrasi yang spesifik; digunakan label konservatif.
- 13: fasilitas kelas ada yang kurang — Teks singkat tentang kekurangan fasilitas tanpa dampak aktivitas yang jelas; digunakan label konservatif.
- 14: kurangnya sistem informasi terkait akademik maupun beberapa pelayanan — Label dipertahankan setelah pemeriksaan konteks teks. Kekurangan informasi akademik disebutkan tanpa contoh dampak; belum cukup untuk urgensi lebih tinggi.
- 23: 1. Terkait nilai konversi 
2. Terkait penggunaan e-learning — Label dipertahankan setelah pemeriksaan konteks teks. Nilai konversi dan e-learning disebut sebagai kendala, tetapi detail dampaknya tidak dijelaskan.
- 26: Sistem kurikulum — Label dipertahankan setelah pemeriksaan konteks teks. Teks hanya menyebut sistem kurikulum tanpa masalah atau dampak yang dapat dinilai.
- 27: Masalah toilet — Teks hanya menyebut masalah toilet tanpa menjelaskan kondisi atau dampaknya.
- 28: Lapangan futsal — Label dipertahankan setelah pemeriksaan konteks teks. Teks hanya menyebut lapangan futsal; belum ada keluhan atau dampak yang jelas.
- 30: Sistem akademik — Teks hanya menyebut sistem akademik tanpa konteks gangguan.
- 33: Masalah organisasi — Label dipertahankan setelah pemeriksaan konteks teks. Teks hanya menyebut masalah organisasi tanpa dampak layanan yang spesifik.
- 38: Ac diruang kelas — Teks hanya menyebut AC di ruang kelas tanpa menjelaskan kerusakan atau dampak.
- 46: Terkait mk kesenian — Teks hanya menyebut mata kuliah kesenian tanpa masalah atau dampak yang jelas.

## Catatan evaluasi

Tidak ada respons kuesioner yang memenuhi bukti Tinggi secara kuat. Karena itu tidak ada label Tinggi yang dipaksakan dan tidak ada data kuesioner yang diubah menjadi Tinggi untuk menyeimbangkan kelas. Deteksi Tinggi pada input baru tetap dijaga oleh guardrail keselamatan dan pola dampak akademik kritis.

Dataset lama tidak lagi digunakan dalam proses AI aktif. Urgensi final ditentukan dengan rule-based yang dapat dijelaskan.

## Evaluasi model


Metode: `stratified_k_fold` dengan 5 fold dan random seed `42`.
Accuracy: `0.4630`.

| Kelas | Precision | Recall | F1-score |
|---|---:|---:|---:|
| Rendah | 0.4615 | 0.6923 | 0.5538 |
| Sedang | 0.4667 | 0.2500 | 0.3256 |
| Tinggi | 0.0000 | 0.0000 | 0.0000 |

Confusion matrix:

```json
{
  "Rendah": {
    "Rendah": 18,
    "Sedang": 8,
    "Tinggi": 0
  },
  "Sedang": {
    "Rendah": 21,
    "Sedang": 7,
    "Tinggi": 0
  },
  "Tinggi": {
    "Rendah": 0,
    "Sedang": 0,
    "Tinggi": 0
  }
}
```

Kelas Tinggi tidak memiliki sampel reviewed, sehingga precision/recall/F1 Tinggi bernilai 0 dan tidak boleh ditafsirkan sebagai evaluasi kemampuan deteksi Tinggi. Guardrail rule-based dievaluasi melalui test kasus keselamatan dan dampak akademik kritis.
