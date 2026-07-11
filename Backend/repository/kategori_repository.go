package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type KategoriRepository interface {
	Create(kategori *models.KategoriPengaduan) error
	GetByID(id uint64) (*models.KategoriPengaduan, error)
	GetByNama(nama string) (*models.KategoriPengaduan, error)
	GetAll() ([]models.KategoriPengaduan, error)
	Update(kategori *models.KategoriPengaduan) error
	Delete(id uint64) error
}

type kategoriRepository struct {
	db *gorm.DB
}

func NewKategoriRepository() KategoriRepository {
	return &kategoriRepository{
		db: config.DB,
	}
}

// =======================
// CREATE
// =======================

func (r *kategoriRepository) Create(kategori *models.KategoriPengaduan) error {
	return r.db.Create(kategori).Error
}

// =======================
// GET BY ID
// =======================

func (r *kategoriRepository) GetByID(id uint64) (*models.KategoriPengaduan, error) {

	var kategori models.KategoriPengaduan

	err := r.db.First(&kategori, id).Error
	if err != nil {
		return nil, err
	}

	return &kategori, nil
}

// =======================
// GET BY NAMA
// =======================

func (r *kategoriRepository) GetByNama(nama string) (*models.KategoriPengaduan, error) {

	var kategori models.KategoriPengaduan

	err := r.db.Where("nama = ?", nama).First(&kategori).Error
	if err != nil {
		return nil, err
	}

	return &kategori, nil
}

// =======================
// GET ALL
// =======================

func (r *kategoriRepository) GetAll() ([]models.KategoriPengaduan, error) {

	var kategori []models.KategoriPengaduan

	err := r.db.Order("nama ASC").Find(&kategori).Error

	return kategori, err
}

// =======================
// UPDATE
// =======================

func (r *kategoriRepository) Update(kategori *models.KategoriPengaduan) error {
	return r.db.Save(kategori).Error
}

// =======================
// DELETE
// =======================

func (r *kategoriRepository) Delete(id uint64) error {
	return r.db.Delete(&models.KategoriPengaduan{}, id).Error
}
