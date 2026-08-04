package models

type KategoriPengaduan struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama      string `gorm:"column:nama;size:100;unique;not null" json:"nama"`
	Deskripsi string `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	IsActive  bool   `gorm:"column:is_active;default:true" json:"is_active"`

	// Relasi
	Pengaduan []Pengaduan `gorm:"foreignKey:KategoriID" json:"-"`
}

func (KategoriPengaduan) TableName() string {
	return "kategori_pengaduan"
}
