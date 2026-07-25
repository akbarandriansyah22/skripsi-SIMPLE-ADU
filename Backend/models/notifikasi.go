package models

import "time"

type Notifikasi struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"column:user_id;not null" json:"user_id"`
	PengaduanID *uint      `gorm:"column:pengaduan_id" json:"pengaduan_id,omitempty"`
	Judul       string     `gorm:"column:judul;size:255;not null" json:"judul"`
	Isi         string     `gorm:"column:isi;type:text;not null" json:"isi"`
	IsRead      bool       `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	ReadAt      *time.Time `gorm:"column:read_at" json:"read_at,omitempty"`

	// Relasi
	User      User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Pengaduan *Pengaduan `gorm:"foreignKey:PengaduanID;references:ID" json:"-"`
}

func (Notifikasi) TableName() string {
	return "notifikasi"
}
