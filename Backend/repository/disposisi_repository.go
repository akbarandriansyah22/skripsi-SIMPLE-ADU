package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type DisposisiRepository interface {
	Create(disposisi *models.Disposisi) error

	GetByID(id uint64) (*models.Disposisi, error)
	GetByPengaduanID(pengaduanID uint64) (*models.Disposisi, error)
	GetByPimpinanID(pimpinanID uint64) ([]models.Disposisi, error)
	GetAll() ([]models.Disposisi, error)

	Update(disposisi *models.Disposisi) error

	Delete(id uint64) error
}

type disposisiRepository struct {
	db *gorm.DB
}

func NewDisposisiRepository() DisposisiRepository {
	return &disposisiRepository{
		db: config.DB,
	}
}

// ===================================
// CREATE
// ===================================

func (r *disposisiRepository) Create(disposisi *models.Disposisi) error {
	return r.db.Create(disposisi).Error
}

// ===================================
// GET BY ID
// ===================================

func (r *disposisiRepository) GetByID(id uint64) (*models.Disposisi, error) {

	var disposisi models.Disposisi

	err := r.db.
		Preload("Pengaduan").
		Preload("Pimpinan").
		First(&disposisi, id).Error

	if err != nil {
		return nil, err
	}

	return &disposisi, nil
}

// ===================================
// GET BY PENGADUAN ID
// ===================================

func (r *disposisiRepository) GetByPengaduanID(pengaduanID uint64) (*models.Disposisi, error) {

	var disposisi models.Disposisi

	err := r.db.
		Where("pengaduan_id = ?", pengaduanID).
		Preload("Pengaduan").
		Preload("Pimpinan").
		First(&disposisi).Error

	if err != nil {
		return nil, err
	}

	return &disposisi, nil
}

// ===================================
// GET BY PIMPINAN ID
// ===================================

func (r *disposisiRepository) GetByPimpinanID(pimpinanID uint64) ([]models.Disposisi, error) {

	var disposisi []models.Disposisi

	err := r.db.
		Where("pimpinan_id = ?", pimpinanID).
		Preload("Pengaduan").
		Preload("Pimpinan").
		Order("created_at DESC").
		Find(&disposisi).Error

	return disposisi, err
}

// ===================================
// GET ALL
// ===================================

func (r *disposisiRepository) GetAll() ([]models.Disposisi, error) {

	var disposisi []models.Disposisi

	err := r.db.
		Preload("Pengaduan").
		Preload("Pimpinan").
		Order("created_at DESC").
		Find(&disposisi).Error

	return disposisi, err
}

// ===================================
// UPDATE
// ===================================

func (r *disposisiRepository) Update(disposisi *models.Disposisi) error {
	return r.db.
		Model(&models.Disposisi{}).
		Where("id = ?", disposisi.ID).
		Updates(disposisi).Error
}

// ===================================
// DELETE
// ===================================

func (r *disposisiRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Disposisi{}, id).Error
}
