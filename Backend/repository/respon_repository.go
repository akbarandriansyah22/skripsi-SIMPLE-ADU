package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type ResponRepository interface {
	Create(respon *models.ResponPengaduan) error

	GetByID(id uint64) (*models.ResponPengaduan, error)
	GetByPengaduanID(pengaduanID uint64) ([]models.ResponPengaduan, error)
	GetByUserID(userID uint64) ([]models.ResponPengaduan, error)
	GetAll() ([]models.ResponPengaduan, error)

	Update(respon *models.ResponPengaduan) error

	Delete(id uint64) error
}

type responRepository struct {
	db *gorm.DB
}

func NewResponRepository() ResponRepository {
	return &responRepository{
		db: config.DB,
	}
}

// ===================================
// CREATE
// ===================================

func (r *responRepository) Create(respon *models.ResponPengaduan) error {
	return r.db.Create(respon).Error
}

// ===================================
// GET BY ID
// ===================================

func (r *responRepository) GetByID(id uint64) (*models.ResponPengaduan, error) {

	var respon models.ResponPengaduan

	err := r.db.
		Preload("Pengaduan").
		Preload("User").
		First(&respon, id).Error

	if err != nil {
		return nil, err
	}

	return &respon, nil
}

// ===================================
// GET BY PENGADUAN ID
// ===================================

func (r *responRepository) GetByPengaduanID(pengaduanID uint64) ([]models.ResponPengaduan, error) {

	var respon []models.ResponPengaduan

	err := r.db.
		Where("pengaduan_id = ?", pengaduanID).
		Preload("User").
		Order("created_at ASC").
		Find(&respon).Error

	return respon, err
}

// ===================================
// GET BY USER ID
// ===================================

func (r *responRepository) GetByUserID(userID uint64) ([]models.ResponPengaduan, error) {

	var respon []models.ResponPengaduan

	err := r.db.
		Where("user_id = ?", userID).
		Preload("Pengaduan").
		Order("created_at DESC").
		Find(&respon).Error

	return respon, err
}

// ===================================
// GET ALL
// ===================================

func (r *responRepository) GetAll() ([]models.ResponPengaduan, error) {

	var respon []models.ResponPengaduan

	err := r.db.
		Preload("Pengaduan").
		Preload("User").
		Order("created_at DESC").
		Find(&respon).Error

	return respon, err
}

// ===================================
// UPDATE
// ===================================

func (r *responRepository) Update(respon *models.ResponPengaduan) error {
	return r.db.
		Model(&models.ResponPengaduan{}).
		Where("id = ?", respon.ID).
		Updates(respon).Error
}

// ===================================
// DELETE
// ===================================

func (r *responRepository) Delete(id uint64) error {
	return r.db.Delete(&models.ResponPengaduan{}, id).Error
}