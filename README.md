# SIMPEL-ADU

## Setup Lokal Docker

```bash
cp .env.example .env
docker compose up --build
```

Untuk menjalankan backend tanpa Docker, gunakan contoh konfigurasi:

```bash
cp Backend/.env.example Backend/.env
```

File `.env` berisi secret lokal dan tidak boleh disimpan ke Git.

## Konfigurasi PostgreSQL host

Untuk DBeaver/PgTools dari host gunakan:

- Host: `127.0.0.1`
- Port: `5434`
- Database: `simpel_adu`
- Username: nilai `POSTGRES_USER` pada `.env`
- Password: nilai aktif PostgreSQL pada `.env`
- SSL mode: `disable`

Di dalam container backend, gunakan `DB_HOST=postgres` dan `DB_PORT=5432`.
Perubahan schema dijalankan melalui migration di `Backend/migrations/`; volume
`simpel_adu_postgres_data` tidak perlu dihapus saat konfigurasi berubah.

## Role dan akun demo

Role tersimpan secara canonical sebagai `mahasiswa`, `admin_sistem`,
`admin_fakultas`, `kasubag`, dan `pimpinan_fakultas`. Dua akun Kasubag demo
terhubung ke unit `Akademik` dan `Sarana dan Prasarana`. Password demo dibaca
dari `SEED_DEMO_PASSWORD` dan tidak disimpan di Git.
