package models

import "time"

type RiwayatStatusPengaduan struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PengaduanID uint      `gorm:"column:pengaduan_id;not null;index" json:"pengaduan_id"`
	ChangedBy   uint      `gorm:"column:changed_by;not null" json:"changed_by"`
	StatusLama  string    `gorm:"column:status_lama;size:30" json:"status_lama"`
	StatusBaru  string    `gorm:"column:status_baru;size:30;not null" json:"status_baru"`
	Catatan     string    `gorm:"column:catatan;type:text" json:"catatan"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`

	Pengaduan Pengaduan `gorm:"foreignKey:PengaduanID;references:ID" json:"-"`
	User      User      `gorm:"foreignKey:ChangedBy;references:ID" json:"user,omitempty"`
}

func (RiwayatStatusPengaduan) TableName() string { return "riwayat_status_pengaduan" }
