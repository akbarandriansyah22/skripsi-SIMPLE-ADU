package models

import "time"

type HasilAI struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID   uint      `gorm:"column:pengaduan_id;unique;not null" json:"pengaduan_id"`
	Sentimen      string    `gorm:"column:sentimen;size:20;not null" json:"sentimen"`
	SkorSentimen  int       `gorm:"column:skor_sentimen;not null" json:"skor_sentimen"`
	Urgensi       string    `gorm:"column:urgensi;size:20;not null" json:"urgensi"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`

	// Relasi
	Pengaduan *Pengaduan `gorm:"foreignKey:PengaduanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (HasilAI) TableName() string {
	return "hasil_ai"
}