package dto

import "time"

type KoordinasiResponse struct {
	ID               uint      `json:"id"`
	PengaduanID      uint      `json:"pengaduan_id"`
	ParentID         *uint     `json:"parent_id,omitempty"`
	SenderID         uint      `json:"sender_id"`
	SenderName       string    `json:"sender_name"`
	SenderRole       string    `json:"sender_role"`
	Pesan            string    `json:"pesan"`
	Lampiran         string    `json:"lampiran,omitempty"`
	LampiranNamaAsli string    `json:"lampiran_nama_asli,omitempty"`
	LampiranURL      string    `json:"lampiran_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
