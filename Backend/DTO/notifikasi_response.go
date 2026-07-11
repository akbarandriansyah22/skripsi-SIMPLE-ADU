package dto

import "time"

type NotifikasiResponse struct {
	ID        uint      `json:"id"`
	Judul     string    `json:"judul"`
	Isi       string    `json:"isi"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
