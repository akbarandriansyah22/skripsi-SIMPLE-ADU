package repository

import (
	"backend/config"
	"backend/models"

	"gorm.io/gorm"
)

type UnitRepository interface {
	Create(unit *models.Unit) error
	GetByID(id uint64) (*models.Unit, error)
	GetByNama(nama string) (*models.Unit, error)
	GetAll() ([]models.Unit, error)
	Update(unit *models.Unit) error
	Delete(id uint64) error
}

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository() UnitRepository {
	return &unitRepository{
		db: config.DB,
	}
}

// =======================
// CREATE
// =======================

func (r *unitRepository) Create(unit *models.Unit) error {
	return r.db.Create(unit).Error
}

// =======================
// GET BY ID
// =======================

func (r *unitRepository) GetByID(id uint64) (*models.Unit, error) {

	var unit models.Unit

	err := r.db.First(&unit, id).Error
	if err != nil {
		return nil, err
	}

	return &unit, nil
}

// =======================
// GET BY NAMA
// =======================

func (r *unitRepository) GetByNama(nama string) (*models.Unit, error) {

	var unit models.Unit

	err := r.db.Where("nama_unit = ?", nama).First(&unit).Error
	if err != nil {
		return nil, err
	}

	return &unit, nil
}

// =======================
// GET ALL
// =======================

func (r *unitRepository) GetAll() ([]models.Unit, error) {

	var units []models.Unit

	err := r.db.Order("nama_unit ASC").Find(&units).Error

	return units, err
}

// =======================
// UPDATE
// =======================

func (r *unitRepository) Update(unit *models.Unit) error {
	return r.db.Save(unit).Error
}

// =======================
// DELETE
// =======================

func (r *unitRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Unit{}, id).Error
}