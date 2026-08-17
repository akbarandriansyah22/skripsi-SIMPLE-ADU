ALTER TABLE users ADD COLUMN IF NOT EXISTS must_change_password boolean NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS sumber_akun varchar(20) NOT NULL DEFAULT 'manual';
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at timestamptz;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'password_must_change'
    ) THEN
        UPDATE users
        SET must_change_password = password_must_change
        WHERE must_change_password IS DISTINCT FROM password_must_change;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS import_mahasiswa (
    id bigserial PRIMARY KEY,
    admin_sistem_id bigint NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    nama_file varchar(255) NOT NULL,
    status varchar(20) NOT NULL DEFAULT 'Diproses',
    total_data integer NOT NULL DEFAULT 0,
    jumlah_berhasil integer NOT NULL DEFAULT 0,
    jumlah_gagal integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamptz,
    CONSTRAINT import_mahasiswa_status_check CHECK (status IN ('Diproses', 'Selesai', 'Gagal')),
    CONSTRAINT import_mahasiswa_completed_check CHECK (
        (status = 'Diproses' AND completed_at IS NULL)
        OR (status IN ('Selesai', 'Gagal') AND completed_at IS NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS detail_import_mahasiswa (
    id bigserial PRIMARY KEY,
    import_id bigint NOT NULL REFERENCES import_mahasiswa(id) ON UPDATE CASCADE ON DELETE CASCADE,
    nomor_baris integer NOT NULL,
    nama_lengkap varchar(150),
    nim varchar(20),
    email varchar(150),
    program_studi varchar(100),
    angkatan integer,
    status varchar(20) NOT NULL,
    user_id bigint REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
    pesan_error text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT detail_import_nomor_baris_check CHECK (nomor_baris >= 2),
    CONSTRAINT detail_import_status_check CHECK (status IN ('Berhasil', 'Gagal')),
    CONSTRAINT detail_import_outcome_check CHECK (
        (status = 'Berhasil' AND user_id IS NOT NULL AND pesan_error IS NULL)
        OR (status = 'Gagal' AND user_id IS NULL AND pesan_error IS NOT NULL)
    )
);
