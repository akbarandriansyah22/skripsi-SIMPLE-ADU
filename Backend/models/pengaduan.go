package models

import "time"

type Pengaduan struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	KodeTiket        string     `gorm:"column:kode_tiket;size:30;unique;not null" json:"kode_tiket"`
	UserID           uint       `gorm:"column:user_id;not null" json:"user_id"`
	KategoriID       uint       `gorm:"column:kategori_id;not null" json:"kategori_id"`
	UnitID           *uint      `gorm:"column:unit_id" json:"unit_id"`
	Judul            string     `gorm:"column:judul;size:255;not null" json:"judul"`
	Deskripsi        string     `gorm:"column:deskripsi;type:text;not null" json:"deskripsi"`
	Lampiran         string     `gorm:"column:lampiran;type:text" json:"lampiran"`
	Status           string     `gorm:"column:status;size:30;default:'Menunggu Verifikasi'" json:"status"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	TanggalSelesai   *time.Time `gorm:"column:tanggal_selesai" json:"tanggal_selesai"`

	// ===========================
	// Relasi
	// ===========================

	User      User                `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Kategori  KategoriPengaduan   `gorm:"foreignKey:KategoriID;references:ID" json:"kategori"`
	Unit      Unit                `gorm:"foreignKey:UnitID;references:ID" json:"unit"`

	HasilAI   HasilAI             `gorm:"foreignKey:PengaduanID;references:ID" json:"hasil_ai"`
	Disposisi *Disposisi          `gorm:"foreignKey:PengaduanID;references:ID" json:"disposisi"`

	ResponPengaduan []ResponPengaduan `gorm:"foreignKey:PengaduanID;references:ID" json:"respon_pengaduan"`
}

func (Pengaduan) TableName() string {
	return "pengaduan"
}