package repository

import (
	"backend/config"
	"backend/models"
	"strings"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		DB: config.DB,
	}
}

// =========================
// GET USER BY EMAIL
// =========================

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.DB.
		Where("LOWER(email) = ?", strings.ToLower(strings.TrimSpace(email))).
		Preload("Unit").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =========================
// GET USER BY EMAIL OR NIM
// =========================

func (r *AuthRepository) GetUserByEmailOrNIM(identifier string) (*models.User, error) {

	var user models.User

	identifier = strings.TrimSpace(identifier)
	if strings.Contains(identifier, "@") {
		identifier = strings.ToLower(identifier)
	}
	err := r.DB.
		Joins("LEFT JOIN mahasiswa ON mahasiswa.user_id = users.id").
		Preload("Unit").
		Where("LOWER(users.email) = ? OR mahasiswa.nim = ?", identifier, identifier).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =========================
// REGISTER MAHASISWA
// =========================

func (r *AuthRepository) RegisterMahasiswa(
	user *models.User,
	mahasiswa *models.Mahasiswa,
) error {

	tx := r.DB.Begin()

	// Simpan user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Ambil user_id
	mahasiswa.UserID = user.ID

	// Simpan data mahasiswa
	if err := tx.Create(mahasiswa).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// =========================
// GET USER BY ID
// =========================

func (r *AuthRepository) GetUserByID(id uint) (*models.User, error) {

	var user models.User

	err := r.DB.
		Where("id = ?", id).
		Preload("Unit").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =========================
// GET MAHASISWA BY NIM
// =========================

func (r *AuthRepository) GetMahasiswaByNIM(nim string) (*models.Mahasiswa, error) {

	var mahasiswa models.Mahasiswa

	err := r.DB.
		Where("nim = ?", nim).
		First(&mahasiswa).Error

	if err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}
