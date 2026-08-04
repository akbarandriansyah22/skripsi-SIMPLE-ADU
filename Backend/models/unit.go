package models

type Unit struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaUnit string `gorm:"column:nama_unit;size:150;unique;not null" json:"nama_unit"`
	Email    string `gorm:"column:email;size:150" json:"email"`
	IsActive bool   `gorm:"column:is_active;default:true" json:"is_active"`

	// Relasi
	Pengaduan []Pengaduan `gorm:"foreignKey:UnitID" json:"-"`
	Users     []User      `gorm:"foreignKey:UnitID" json:"-"`
}

func (Unit) TableName() string {
	return "unit"
}
