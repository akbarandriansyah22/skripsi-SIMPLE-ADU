package models

import "time"

type ResponPengaduan struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID  uint      `gorm:"column:pengaduan_id;not null" json:"pengaduan_id"`
	UserID       uint      `gorm:"column:user_id;not null" json:"user_id"`
	Pesan        string    `gorm:"column:pesan;type:text;not null" json:"pesan"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`

	// Relasi
	Pengaduan Pengaduan `gorm:"foreignKey:PengaduanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"pengaduan"`
	User       User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}

func (ResponPengaduan) TableName() string {
	return "respon_pengaduan"
}