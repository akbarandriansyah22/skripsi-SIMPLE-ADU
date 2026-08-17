package models

type Mahasiswa struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint   `gorm:"column:user_id;unique;not null" json:"user_id"`
	NIM          string `gorm:"column:nim;size:20;unique;not null" json:"nim"`
	ProgramStudi string `gorm:"column:program_studi;size:100;not null" json:"program_studi"`
	Angkatan     int    `gorm:"column:angkatan;not null" json:"angkatan"`
	NoHP         *string `gorm:"column:no_hp;size:20" json:"no_hp,omitempty"`

	// Relasi
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}

func (Mahasiswa) TableName() string {
	return "mahasiswa"
}
