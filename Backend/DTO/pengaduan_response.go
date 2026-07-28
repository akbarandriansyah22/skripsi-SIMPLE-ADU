package dto

import "time"

type PengaduanResponse struct {
	ID               uint                      `json:"id"`
	KodeTiket        string                    `json:"kode_tiket"`
	UserID           uint                      `json:"user_id"`
	KategoriID       uint                      `json:"kategori_id"`
	Kategori         *KategoriResponse         `json:"kategori,omitempty"`
	UnitID           *uint                     `json:"unit_id,omitempty"`
	Unit             *UnitResponse             `json:"unit,omitempty"`
	User             *UserResponse             `json:"user,omitempty"`
	Judul            string                    `json:"judul"`
	Deskripsi        string                    `json:"deskripsi"`
	Lampiran         string                    `json:"lampiran"`
	LampiranNamaAsli string                    `json:"lampiran_nama_asli,omitempty"`
	LampiranMimeType string                    `json:"lampiran_mime_type,omitempty"`
	LampiranUkuran   int64                     `json:"lampiran_ukuran,omitempty"`
	LampiranURL      string                    `json:"lampiran_url,omitempty"`
	Status           string                    `json:"status"`
	SkorSentimen     *int                      `json:"skor_sentimen,omitempty"`
	Sentimen         string                    `json:"sentimen,omitempty"`
	Urgensi          string                    `json:"urgensi,omitempty"`
	AIStatus         string                    `json:"ai_status"`
	Validasi         *ValidasiResponse         `json:"validasi,omitempty"`
	Disposisi        *DisposisiResponse        `json:"disposisi,omitempty"`
	RiwayatStatus    []RiwayatStatusResponse   `json:"riwayat_status_pengaduan,omitempty"`
	ResponPengaduan  []ResponPengaduanResponse `json:"respon_pengaduan,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
}

type KategoriResponse struct {
	ID   uint   `json:"id"`
	Nama string `json:"nama"`
}

type UnitResponse struct {
	ID       uint   `json:"id"`
	NamaUnit string `json:"nama_unit"`
}

type ValidasiResponse struct {
	ID              uint   `json:"id"`
	AdminFakultasID uint   `json:"admin_fakultas_id"`
	StatusValidasi  string `json:"status_validasi"`
	Catatan         string `json:"catatan"`
}

type DisposisiResponse struct {
	ID         uint   `json:"id"`
	PimpinanID uint   `json:"pimpinan_id"`
	UnitID     uint   `json:"unit_id"`
	Catatan    string `json:"catatan"`
}

type RiwayatStatusResponse struct {
	ID         uint      `json:"id"`
	ChangedBy  uint      `json:"changed_by"`
	StatusLama string    `json:"status_lama"`
	StatusBaru string    `json:"status_baru"`
	Catatan    string    `json:"catatan"`
	CreatedAt  time.Time `json:"created_at"`
}

type ReanalyzeAIResponse struct {
	PengaduanID  uint   `json:"pengaduan_id"`
	SkorSentimen int    `json:"skor_sentimen"`
	Sentimen     string `json:"sentimen"`
	Urgensi      string `json:"urgensi"`
	AIStatus     string `json:"ai_status"`
}

type ResponPengaduanResponse struct {
	ID               uint          `json:"id"`
	PengaduanID      uint          `json:"pengaduan_id"`
	UserID           uint          `json:"user_id"`
	Pesan            string        `json:"pesan"`
	Lampiran         string        `json:"lampiran,omitempty"`
	LampiranNamaAsli string        `json:"lampiran_nama_asli,omitempty"`
	LampiranMimeType string        `json:"lampiran_mime_type,omitempty"`
	LampiranUkuran   int64         `json:"lampiran_ukuran,omitempty"`
	LampiranURL      string        `json:"lampiran_url,omitempty"`
	User             *UserResponse `json:"user,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
}
