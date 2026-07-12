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
