package models

import "time"

// KoordinasiInternal menyimpan percakapan privat antara pimpinan dan unit.
// Pesan ini tidak pernah dimasukkan ke ResponPengaduan yang terlihat mahasiswa.
type KoordinasiInternal struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID      uint      `gorm:"column:pengaduan_id;not null;index" json:"pengaduan_id"`
	SenderID         uint      `gorm:"column:sender_id;not null;index" json:"sender_id"`
	ParentID         *uint     `gorm:"column:parent_id;index" json:"parent_id,omitempty"`
	Pesan            string    `gorm:"column:pesan;type:text;not null" json:"pesan"`
	Lampiran         string    `gorm:"column:lampiran;type:text" json:"lampiran,omitempty"`
	LampiranNamaAsli *string   `gorm:"column:lampiran_nama_asli;type:text" json:"lampiran_nama_asli,omitempty"`
	LampiranMimeType *string   `gorm:"column:lampiran_mime_type;size:100" json:"lampiran_mime_type,omitempty"`
	LampiranUkuran   *int64    `gorm:"column:lampiran_ukuran" json:"lampiran_ukuran,omitempty"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`

	Pengaduan Pengaduan `gorm:"foreignKey:PengaduanID;references:ID" json:"-"`
	Sender    User      `gorm:"foreignKey:SenderID;references:ID" json:"sender,omitempty"`
}

func (KoordinasiInternal) TableName() string { return "koordinasi_internal" }
