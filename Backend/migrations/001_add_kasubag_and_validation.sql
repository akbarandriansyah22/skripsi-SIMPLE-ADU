CREATE TABLE IF NOT EXISTS schema_migrations (
    version varchar(255) PRIMARY KEY,
    applied_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS unit_id bigint;
CREATE INDEX IF NOT EXISTS idx_users_unit_id ON users (unit_id);

ALTER TABLE disposisi ADD COLUMN IF NOT EXISTS unit_id bigint;
CREATE INDEX IF NOT EXISTS idx_disposisi_unit_id ON disposisi (unit_id);

CREATE TABLE IF NOT EXISTS validasi_pengaduan (
    id bigserial PRIMARY KEY,
    pengaduan_id bigint NOT NULL UNIQUE,
    admin_fakultas_id bigint NOT NULL,
    status_validasi varchar(20) NOT NULL,
    catatan text,
    validated_at timestamptz,
    updated_at timestamptz
);
ALTER TABLE validasi_pengaduan ADD COLUMN IF NOT EXISTS admin_fakultas_id bigint;
ALTER TABLE validasi_pengaduan ADD COLUMN IF NOT EXISTS validated_at timestamptz;
ALTER TABLE validasi_pengaduan ADD COLUMN IF NOT EXISTS updated_at timestamptz;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'validasi_pengaduan' AND column_name = 'admin_id') THEN
        UPDATE validasi_pengaduan SET admin_fakultas_id = admin_id WHERE admin_fakultas_id IS NULL;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'validasi_pengaduan' AND column_name = 'created_at') THEN
        UPDATE validasi_pengaduan SET validated_at = created_at WHERE validated_at IS NULL;
        UPDATE validasi_pengaduan SET updated_at = COALESCE(updated_at, created_at, CURRENT_TIMESTAMP) WHERE updated_at IS NULL;
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_validasi_pengaduan_pengaduan_id ON validasi_pengaduan (pengaduan_id);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'users'::regclass AND confrelid = 'unit'::regclass AND contype = 'f') THEN
        ALTER TABLE users ADD CONSTRAINT fk_users_unit FOREIGN KEY (unit_id) REFERENCES unit(id) ON UPDATE CASCADE ON DELETE SET NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'disposisi'::regclass AND confrelid = 'unit'::regclass AND contype = 'f') THEN
        ALTER TABLE disposisi ADD CONSTRAINT fk_disposisi_unit FOREIGN KEY (unit_id) REFERENCES unit(id) ON UPDATE CASCADE ON DELETE SET NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'validasi_pengaduan'::regclass AND confrelid = 'pengaduan'::regclass AND contype = 'f') THEN
        ALTER TABLE validasi_pengaduan ADD CONSTRAINT fk_validasi_pengaduan_pengaduan FOREIGN KEY (pengaduan_id) REFERENCES pengaduan(id) ON UPDATE CASCADE ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'validasi_pengaduan'::regclass AND conname = 'fk_validasi_pengaduan_admin_fakultas') THEN
        ALTER TABLE validasi_pengaduan ADD CONSTRAINT fk_validasi_pengaduan_admin_fakultas FOREIGN KEY (admin_fakultas_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;
    END IF;
END $$;
