package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: config.DB,
	}
}

// =========================
// CREATE USER
// =========================

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// =========================
// GET USER BY ID
// =========================

func (r *userRepository) GetByID(id uint64) (*models.User, error) {

	var user models.User

	err := r.db.First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =========================
// GET USER BY EMAIL
// =========================

func (r *userRepository) GetByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =========================
// GET ALL USER
// =========================

func (r *userRepository) GetAll() ([]models.User, error) {

	var users []models.User

	err := r.db.Find(&users).Error

	return users, err
}

// =========================
// UPDATE USER
// =========================

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// =========================
// DELETE USER
// =========================

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}
