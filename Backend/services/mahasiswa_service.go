package services

import (
	dto "backend/DTO"
	"backend/models"
	"backend/repository"
)

type MahasiswaService struct {
	repo repository.MahasiswaRepository
}

func NewMahasiswaService() *MahasiswaService {
	return &MahasiswaService{
		repo: repository.NewMahasiswaRepository(),
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

	mahasiswa, err := s.repo.GetByUserID(userID)

	if err != nil {
		return err
	}

	mahasiswa.NoHP = req.NoHP

	return s.repo.Update(mahasiswa)
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