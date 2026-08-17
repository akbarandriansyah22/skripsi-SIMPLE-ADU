package services

import (
	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/repository"

	"gorm.io/gorm"
)

type MahasiswaService struct {
	repo     repository.MahasiswaRepository
	userRepo repository.UserRepository
}

func NewMahasiswaService() *MahasiswaService {
	return &MahasiswaService{
		repo:     repository.NewMahasiswaRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// ======================================
// PROFILE
// ======================================

func (s *MahasiswaService) GetProfile(userID uint64) (*models.Mahasiswa, error) {
	return s.repo.GetByUserID(userID)
}

// ======================================
// UPDATE PROFILE
// ======================================

func (s *MahasiswaService) UpdateProfile(userID uint64, req dto.UpdateProfileRequest) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var mahasiswa models.Mahasiswa
		if err := tx.Where("user_id = ?", userID).First(&mahasiswa).Error; err != nil {
			return err
		}

		mahasiswa.NoHP = optionalString(req.NoHP)
		if err := tx.Save(&mahasiswa).Error; err != nil {
			return err
		}

		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		user.NamaLengkap = req.NamaLengkap
		return tx.Save(&user).Error
	})
}

// ======================================
// GET BY NIM
// ======================================

func (s *MahasiswaService) GetByNIM(nim string) (*models.Mahasiswa, error) {
	return s.repo.GetByNIM(nim)
}

// ======================================
// GET ALL
// ======================================

func (s *MahasiswaService) GetAll() ([]models.Mahasiswa, error) {
	return s.repo.GetAll()
}

// ======================================
// DELETE
// ======================================

func (s *MahasiswaService) Delete(id uint64) error {
	return s.repo.Delete(id)
}
