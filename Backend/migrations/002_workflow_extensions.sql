ALTER TABLE users ADD COLUMN IF NOT EXISTS must_change_password boolean NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS unit_id bigint;
ALTER TABLE kategori_pengaduan ADD COLUMN IF NOT EXISTS is_active boolean NOT NULL DEFAULT true;
ALTER TABLE unit ADD COLUMN IF NOT EXISTS is_active boolean NOT NULL DEFAULT true;
ALTER TABLE notifikasi ADD COLUMN IF NOT EXISTS read_at timestamptz;
ALTER TABLE pengaduan ADD COLUMN IF NOT EXISTS lampiran_nama_asli text;
ALTER TABLE pengaduan ADD COLUMN IF NOT EXISTS lampiran_mime_type varchar(100);
ALTER TABLE pengaduan ADD COLUMN IF NOT EXISTS lampiran_ukuran bigint;
ALTER TABLE respon_pengaduan ADD COLUMN IF NOT EXISTS lampiran_nama_asli text;
ALTER TABLE respon_pengaduan ADD COLUMN IF NOT EXISTS lampiran_mime_type varchar(100);
ALTER TABLE respon_pengaduan ADD COLUMN IF NOT EXISTS lampiran_ukuran bigint;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS cleaned_text text;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS tokens jsonb;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS skor_positif bigint NOT NULL DEFAULT 0;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS skor_negatif bigint NOT NULL DEFAULT 0;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS penjelasan_sentimen text;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS detail_skor jsonb;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS matched_words jsonb;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS urgency_score bigint NOT NULL DEFAULT 0;
ALTER TABLE hasil_ai ADD COLUMN IF NOT EXISTS urgency_reason text;

CREATE TABLE IF NOT EXISTS riwayat_status_pengaduan (
    id bigserial PRIMARY KEY,
    pengaduan_id bigint NOT NULL REFERENCES pengaduan(id) ON UPDATE CASCADE ON DELETE CASCADE,
    changed_by bigint NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    status_lama varchar(30),
    status_baru varchar(30) NOT NULL,
    catatan text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_riwayat_status_pengaduan_id ON riwayat_status_pengaduan(pengaduan_id, created_at);

-- Preserve older records that had already passed Admin Fakultas validation
-- before the status timeline table was introduced.  The validation timestamp
-- is the closest authoritative timestamp for their entry into Pimpinan.
INSERT INTO riwayat_status_pengaduan
    (pengaduan_id, changed_by, status_lama, status_baru, catatan, created_at)
SELECT
    p.id,
    v.admin_fakultas_id,
    'Menunggu Verifikasi',
    'Menunggu Disposisi',
    'Riwayat Pimpinan dipulihkan dari data workflow sebelumnya',
    COALESCE(v.validated_at, p.updated_at, p.created_at)
FROM pengaduan p
JOIN hasil_ai h ON h.pengaduan_id = p.id
JOIN validasi_pengaduan v ON v.pengaduan_id = p.id
WHERE LOWER(h.urgensi) = 'tinggi'
  AND LOWER(v.status_validasi) = 'diterima'
  AND v.admin_fakultas_id IS NOT NULL
  AND LOWER(p.status) IN (
      'menunggu disposisi', 'diteruskan ke unit', 'diproses', 'selesai', 'ditolak'
  )
  AND NOT EXISTS (
      SELECT 1
      FROM riwayat_status_pengaduan existing
      WHERE existing.pengaduan_id = p.id
        AND LOWER(existing.status_baru) = 'menunggu disposisi'
  );

CREATE TABLE IF NOT EXISTS koordinasi_internal (
    id bigserial PRIMARY KEY,
    pengaduan_id bigint NOT NULL REFERENCES pengaduan(id) ON UPDATE CASCADE ON DELETE CASCADE,
    sender_id bigint NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    parent_id bigint REFERENCES koordinasi_internal(id) ON UPDATE CASCADE ON DELETE SET NULL,
    pesan text NOT NULL,
    lampiran text,
    lampiran_nama_asli text,
    lampiran_mime_type varchar(100),
    lampiran_ukuran bigint,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_koordinasi_internal_pengaduan ON koordinasi_internal(pengaduan_id, created_at);


