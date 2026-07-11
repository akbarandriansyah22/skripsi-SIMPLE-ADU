package dto

import "time"

type PengaduanResponse struct {
	ID               uint      `json:"id"`
	KodeTiket        string    `json:"kode_tiket"`
	UserID           uint      `json:"user_id"`
	KategoriID       uint      `json:"kategori_id"`
	UnitID           *uint     `json:"unit_id,omitempty"`
	Judul            string    `json:"judul"`
	Deskripsi        string    `json:"deskripsi"`
	Lampiran         string    `json:"lampiran"`
	Status           string    `json:"status"`
	KategoriPrediksi string    `json:"kategori_prediksi,omitempty"`
	Sentimen         string    `json:"sentimen,omitempty"`
	Urgensi          string    `json:"urgensi,omitempty"`
	Confidence       float64   `json:"confidence,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
