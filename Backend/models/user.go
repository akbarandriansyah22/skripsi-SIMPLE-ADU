package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaLengkap  string    `gorm:"column:nama_lengkap;size:150;not null" json:"nama_lengkap"`
	Email        string    `gorm:"column:email;size:150;unique;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null" json:"-"`
	Role         string    `gorm:"column:role;size:20;not null" json:"role"`
	IsActive     bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relasi
	Pengaduan  []Pengaduan  `gorm:"foreignKey:UserID" json:"-"`
	Notifikasi []Notifikasi `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}
