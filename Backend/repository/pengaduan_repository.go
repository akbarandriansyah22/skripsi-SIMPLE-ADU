package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type PengaduanRepository interface {
	Create(pengaduan *models.Pengaduan) error

	GetByID(id uint64) (*models.Pengaduan, error)
	GetByKodeTiket(kode string) (*models.Pengaduan, error)
	GetByUserID(userID uint64) ([]models.Pengaduan, error)
	GetByStatus(status string) ([]models.Pengaduan, error)
	GetAll() ([]models.Pengaduan, error)

	Update(pengaduan *models.Pengaduan) error
	UpdateStatus(id uint64, status string) error

	Delete(id uint64) error
}

type pengaduanRepository struct {
	db *gorm.DB
}

func NewPengaduanRepository() PengaduanRepository {
	return &pengaduanRepository{
		db: config.DB,
	}
}

// ===================================
// CREATE
// ===================================

func (r *pengaduanRepository) Create(pengaduan *models.Pengaduan) error {
	return r.db.Create(pengaduan).Error
}

// ===================================
// GET BY ID
// ===================================

func (r *pengaduanRepository) GetByID(id uint64) (*models.Pengaduan, error) {

	var pengaduan models.Pengaduan

	err := r.db.
		Preload("User").
		Preload("Kategori").
		Preload("Unit").
		Preload("HasilAI").
		Preload("Disposisi").
		Preload("Disposisi.Unit").
		Preload("Validasi").
		Preload("RiwayatStatus.User").
		Preload("ResponPengaduan.User").
		First(&pengaduan, id).Error

	if err != nil {
		return nil, err
	}

	return &pengaduan, nil
}

// ===================================
// GET BY KODE TIKET
// ===================================

func (r *pengaduanRepository) GetByKodeTiket(kode string) (*models.Pengaduan, error) {

	var pengaduan models.Pengaduan

	err := r.db.
		Where("kode_tiket = ?", kode).
		Preload("User").
		Preload("Kategori").
		Preload("Unit").
		Preload("HasilAI").
		Preload("Disposisi").
		Preload("Disposisi.Unit").
		Preload("Validasi").
		Preload("RiwayatStatus.User").
		Preload("ResponPengaduan.User").
		First(&pengaduan).Error

	if err != nil {
		return nil, err
	}

	return &pengaduan, nil
}

// ===================================
// GET BY USER ID
// ===================================

func (r *pengaduanRepository) GetByUserID(userID uint64) ([]models.Pengaduan, error) {

	var pengaduan []models.Pengaduan

	err := r.db.
		Where("user_id = ?", userID).
		Preload("Kategori").
		Preload("Unit").
		Preload("HasilAI").
		Preload("Disposisi").
		Preload("Disposisi.Unit").
		Preload("Validasi").
		Preload("RiwayatStatus.User").
		Order("created_at DESC").
		Find(&pengaduan).Error

	return pengaduan, err
}

// ===================================
// GET BY STATUS
// ===================================

func (r *pengaduanRepository) GetByStatus(status string) ([]models.Pengaduan, error) {

	var pengaduan []models.Pengaduan

	err := r.db.
		Where("status = ?", status).
		Preload("User").
		Preload("Kategori").
		Preload("Unit").
		Preload("HasilAI").
		Preload("Disposisi").
		Preload("Disposisi.Unit").
		Preload("Validasi").
		Preload("RiwayatStatus.User").
		Order("created_at DESC").
		Find(&pengaduan).Error

	return pengaduan, err
}

// ===================================
// GET ALL
// ===================================

func (r *pengaduanRepository) GetAll() ([]models.Pengaduan, error) {

	var pengaduan []models.Pengaduan

	err := r.db.
		Preload("User").
		Preload("Kategori").
		Preload("Unit").
		Preload("HasilAI").
		Preload("Disposisi").
		Preload("Disposisi.Unit").
		Preload("Validasi").
		Preload("RiwayatStatus.User").
		Order("created_at DESC").
		Find(&pengaduan).Error

	return pengaduan, err
}

// ===================================
// UPDATE
// ===================================

func (r *pengaduanRepository) Update(pengaduan *models.Pengaduan) error {
	return r.db.
		Model(&models.Pengaduan{}).
		Where("id = ?", pengaduan.ID).
		Updates(pengaduan).Error
}

// ===================================
// UPDATE STATUS
// ===================================

func (r *pengaduanRepository) UpdateStatus(id uint64, status string) error {
	return r.db.
		Model(&models.Pengaduan{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ===================================
// DELETE
// ===================================

func (r *pengaduanRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Pengaduan{}, id).Error
}
