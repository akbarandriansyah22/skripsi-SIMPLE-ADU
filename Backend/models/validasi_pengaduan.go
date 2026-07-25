package models

import "time"

type ValidasiPengaduan struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID     uint      `gorm:"column:pengaduan_id;unique;not null" json:"pengaduan_id"`
	AdminFakultasID uint      `gorm:"column:admin_fakultas_id;not null" json:"admin_fakultas_id"`
	StatusValidasi  string    `gorm:"column:status_validasi;size:20;not null" json:"status_validasi"`
	Catatan         string    `gorm:"column:catatan;type:text" json:"catatan"`
	ValidatedAt     time.Time `gorm:"column:validated_at" json:"validated_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`

	Pengaduan     Pengaduan `gorm:"foreignKey:PengaduanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	AdminFakultas User      `gorm:"foreignKey:AdminFakultasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"admin_fakultas,omitempty"`
}

func (ValidasiPengaduan) TableName() string { return "validasi_pengaduan" }
