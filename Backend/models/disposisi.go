package models

import "time"

type Disposisi struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID uint      `gorm:"column:pengaduan_id;unique;not null" json:"pengaduan_id"`
	PimpinanID  uint      `gorm:"column:pimpinan_id;not null" json:"pimpinan_id"`
	Catatan     string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`

	// Relasi
	Pengaduan Pengaduan `gorm:"foreignKey:PengaduanID;references:ID" json:"pengaduan"`
	Pimpinan  User       `gorm:"foreignKey:PimpinanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"pimpinan"`
}

func (Disposisi) TableName() string {
	return "disposisi"
}