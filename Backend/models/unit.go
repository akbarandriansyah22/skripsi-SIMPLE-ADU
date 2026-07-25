package models

type Unit struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaUnit string `gorm:"column:nama_unit;size:150;unique;not null" json:"nama_unit"`
	Email    string `gorm:"column:email;size:150" json:"email"`

	// Relasi
	Pengaduan []Pengaduan `gorm:"foreignKey:UnitID" json:"-"`
	Users     []User      `gorm:"foreignKey:UnitID" json:"-"`
}

func (Unit) TableName() string {
	return "unit"
}
