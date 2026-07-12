package dto

import "time"

type PengaduanResponse struct {
	ID               uint                      `json:"id"`
	KodeTiket        string                    `json:"kode_tiket"`
	UserID           uint                      `json:"user_id"`
	KategoriID       uint                      `json:"kategori_id"`
	UnitID           *uint                     `json:"unit_id,omitempty"`
	User             *UserResponse             `json:"user,omitempty"`
	Judul            string                    `json:"judul"`
	Deskripsi        string                    `json:"deskripsi"`
	Lampiran         string                    `json:"lampiran"`
	Status           string                    `json:"status"`
	KategoriPrediksi string                    `json:"kategori_prediksi,omitempty"`
	Sentimen         string                    `json:"sentimen,omitempty"`
	Urgensi          string                    `json:"urgensi,omitempty"`
	Confidence       float64                   `json:"confidence,omitempty"`
	ResponPengaduan  []ResponPengaduanResponse `json:"respon_pengaduan,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
}

type ResponPengaduanResponse struct {
	ID          uint          `json:"id"`
	PengaduanID uint          `json:"pengaduan_id"`
	UserID      uint          `json:"user_id"`
	Pesan       string        `json:"pesan"`
	User        *UserResponse `json:"user,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}
