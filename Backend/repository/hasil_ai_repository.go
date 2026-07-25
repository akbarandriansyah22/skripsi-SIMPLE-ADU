package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HasilAIRepository interface {
	Create(hasil *models.HasilAI) error
	UpsertByPengaduanID(hasil *models.HasilAI) error

	GetByID(id uint64) (*models.HasilAI, error)
	GetByPengaduanID(pengaduanID uint64) (*models.HasilAI, error)
	GetAll() ([]models.HasilAI, error)

	Update(hasil *models.HasilAI) error

	Delete(id uint64) error
}

type hasilAIRepository struct {
	db *gorm.DB
}

func NewHasilAIRepository() HasilAIRepository {
	return &hasilAIRepository{
		db: config.DB,
	}
}

// ===================================
// CREATE
// ===================================

func (r *hasilAIRepository) Create(hasil *models.HasilAI) error {
	return r.db.Create(hasil).Error
}

func (r *hasilAIRepository) UpsertByPengaduanID(hasil *models.HasilAI) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "pengaduan_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"cleaned_text",
			"tokens",
			"skor_sentimen",
			"sentimen",
			"penjelasan_sentimen",
			"detail_skor",
			"urgensi",
			"dasar_urgensi",
			"updated_at",
		}),
	}).Create(hasil).Error
}

// ===================================
// GET BY ID
// ===================================

func (r *hasilAIRepository) GetByID(id uint64) (*models.HasilAI, error) {

	var hasil models.HasilAI

	err := r.db.
		Preload("Pengaduan").
		First(&hasil, id).Error

	if err != nil {
		return nil, err
	}

	return &hasil, nil
}

// ===================================
// GET BY PENGADUAN ID
// ===================================

func (r *hasilAIRepository) GetByPengaduanID(pengaduanID uint64) (*models.HasilAI, error) {

	var hasil models.HasilAI

	err := r.db.
		Where("pengaduan_id = ?", pengaduanID).
		Preload("Pengaduan").
		First(&hasil).Error

	if err != nil {
		return nil, err
	}

	return &hasil, nil
}

// ===================================
// GET ALL
// ===================================

func (r *hasilAIRepository) GetAll() ([]models.HasilAI, error) {

	var hasil []models.HasilAI

	err := r.db.
		Preload("Pengaduan").
		Order("created_at DESC").
		Find(&hasil).Error

	return hasil, err
}

// ===================================
// UPDATE
// ===================================

func (r *hasilAIRepository) Update(hasil *models.HasilAI) error {
	return r.db.
		Model(&models.HasilAI{}).
		Where("id = ?", hasil.ID).
		Updates(hasil).Error
}

// ===================================
// DELETE
// ===================================

func (r *hasilAIRepository) Delete(id uint64) error {
	return r.db.Delete(&models.HasilAI{}, id).Error
}
