package repository

import (
	"backend/config"
	"backend/models"
	"time"

	"gorm.io/gorm"
)

type NotifikasiRepository interface {
	Create(notifikasi *models.Notifikasi) error

	GetByID(id uint64) (*models.Notifikasi, error)
	GetByUserID(userID uint64) ([]models.Notifikasi, error)
	GetUnreadByUserID(userID uint64) ([]models.Notifikasi, error)
	GetAll() ([]models.Notifikasi, error)

	Update(notifikasi *models.Notifikasi) error
	MarkAsRead(id uint64) error
	MarkAsReadByUserID(id uint64, userID uint64) error

	Delete(id uint64) error
}

type notifikasiRepository struct {
	db *gorm.DB
}

func NewNotifikasiRepository() NotifikasiRepository {
	return &notifikasiRepository{
		db: config.DB,
	}
}

// ===================================
// CREATE
// ===================================

func (r *notifikasiRepository) Create(notifikasi *models.Notifikasi) error {
	return r.db.Create(notifikasi).Error
}

// ===================================
// GET BY ID
// ===================================

func (r *notifikasiRepository) GetByID(id uint64) (*models.Notifikasi, error) {

	var notifikasi models.Notifikasi

	err := r.db.
		Preload("User").
		First(&notifikasi, id).Error

	if err != nil {
		return nil, err
	}

	return &notifikasi, nil
}

// ===================================
// GET BY USER ID
// ===================================

func (r *notifikasiRepository) GetByUserID(userID uint64) ([]models.Notifikasi, error) {

	var notifikasi []models.Notifikasi

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifikasi).Error

	return notifikasi, err
}

// ===================================
// GET UNREAD
// ===================================

func (r *notifikasiRepository) GetUnreadByUserID(userID uint64) ([]models.Notifikasi, error) {

	var notifikasi []models.Notifikasi

	err := r.db.
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifikasi).Error

	return notifikasi, err
}

// ===================================
// GET ALL
// ===================================

func (r *notifikasiRepository) GetAll() ([]models.Notifikasi, error) {

	var notifikasi []models.Notifikasi

	err := r.db.
		Preload("User").
		Order("created_at DESC").
		Find(&notifikasi).Error

	return notifikasi, err
}

// ===================================
// UPDATE
// ===================================

func (r *notifikasiRepository) Update(notifikasi *models.Notifikasi) error {
	return r.db.
		Model(&models.Notifikasi{}).
		Where("id = ?", notifikasi.ID).
		Updates(notifikasi).Error
}

// ===================================
// MARK AS READ
// ===================================

func (r *notifikasiRepository) MarkAsRead(id uint64) error {
	return r.db.
		Model(&models.Notifikasi{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"is_read": true, "read_at": time.Now()}).Error
}

func (r *notifikasiRepository) MarkAsReadByUserID(id uint64, userID uint64) error {
	result := r.db.
		Model(&models.Notifikasi{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ===================================
// DELETE
// ===================================

func (r *notifikasiRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Notifikasi{}, id).Error
}
