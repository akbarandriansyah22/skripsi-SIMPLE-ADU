package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type MahasiswaRepository interface {
	Create(mahasiswa *models.Mahasiswa) error
	GetByID(id uint64) (*models.Mahasiswa, error)
	GetByUserID(userID uint64) (*models.Mahasiswa, error)
	GetByNIM(nim string) (*models.Mahasiswa, error)
	GetAll() ([]models.Mahasiswa, error)
	Update(mahasiswa *models.Mahasiswa) error
	Delete(id uint64) error
}

type mahasiswaRepository struct {
	db *gorm.DB
}

func NewMahasiswaRepository() MahasiswaRepository {
	return &mahasiswaRepository{
		db: config.DB,
	}
}

// =======================
// CREATE
// =======================

func (r *mahasiswaRepository) Create(mahasiswa *models.Mahasiswa) error {
	return r.db.Create(mahasiswa).Error
}

// =======================
// GET BY ID
// =======================

func (r *mahasiswaRepository) GetByID(id uint64) (*models.Mahasiswa, error) {

	var mahasiswa models.Mahasiswa

	err := r.db.First(&mahasiswa, id).Error
	if err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}

// =======================
// GET BY USER ID
// =======================

func (r *mahasiswaRepository) GetByUserID(userID uint64) (*models.Mahasiswa, error) {

	var mahasiswa models.Mahasiswa

	err := r.db.Where("user_id = ?", userID).First(&mahasiswa).Error
	if err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}

// =======================
// GET BY NIM
// =======================

func (r *mahasiswaRepository) GetByNIM(nim string) (*models.Mahasiswa, error) {

	var mahasiswa models.Mahasiswa

	err := r.db.Where("nim = ?", nim).First(&mahasiswa).Error
	if err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}

// =======================
// GET ALL
// =======================

func (r *mahasiswaRepository) GetAll() ([]models.Mahasiswa, error) {

	var mahasiswa []models.Mahasiswa

	err := r.db.Find(&mahasiswa).Error

	return mahasiswa, err
}

// =======================
// UPDATE
// =======================

func (r *mahasiswaRepository) Update(mahasiswa *models.Mahasiswa) error {
	return r.db.Save(mahasiswa).Error
}

// =======================
// DELETE
// =======================

func (r *mahasiswaRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Mahasiswa{}, id).Error
}