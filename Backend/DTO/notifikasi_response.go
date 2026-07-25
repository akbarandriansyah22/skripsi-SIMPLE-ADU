package dto

import "time"

type NotifikasiResponse struct {
	ID          uint       `json:"id"`
	PengaduanID *uint      `json:"pengaduan_id,omitempty"`
	Judul       string     `json:"judul"`
	Isi         string     `json:"isi"`
	IsRead      bool       `json:"is_read"`
	CreatedAt   time.Time  `json:"created_at"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
}
