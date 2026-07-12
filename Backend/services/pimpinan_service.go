package services

import (
	"strings"

	dto "backend/DTO"
	"backend/models"
	"backend/repository"
)

type PimpinanService struct {
	pengaduanRepo repository.PengaduanRepository
	disposisiRepo repository.DisposisiRepository
}

func NewPimpinanService() *PimpinanService {
	return &PimpinanService{
		pengaduanRepo: repository.NewPengaduanRepository(),
		disposisiRepo: repository.NewDisposisiRepository(),
	}
}

func (s *PimpinanService) Dashboard() (*dto.DashboardPimpinanResponse, error) {
	items, err := s.urgensiTinggi()
	if err != nil {
		return nil, err
	}

	resp := &dto.DashboardPimpinanResponse{
		TotalUrgensiTinggi: int64(len(items)),
	}

	for _, item := range items {
		if item.Disposisi == nil {
			resp.BelumDisposisi++
		} else {
			resp.SudahDisposisi++
		}
	}

	return resp, nil
}

func (s *PimpinanService) GetUrgensiTinggi() ([]dto.PengaduanResponse, error) {
	items, err := s.urgensiTinggi()
	if err != nil {
		return nil, err
	}

	return mapPengaduanResponses(items), nil
}

func (s *PimpinanService) CreateDisposisi(pimpinanID uint, pengaduanID uint64, req dto.DisposisiRequest) error {
	pengaduan, err := s.pengaduanRepo.GetByID(pengaduanID)
	if err != nil {
		return err
	}

	disposisi := &models.Disposisi{
		PengaduanID: pengaduan.ID,
		PimpinanID:  pimpinanID,
		Catatan:     req.Catatan,
	}

	if err := s.disposisiRepo.Create(disposisi); err != nil {
		return err
	}

	return s.pengaduanRepo.UpdateStatus(pengaduanID, "Diproses")
}

func (s *PimpinanService) GetDisposisiByPimpinan(pimpinanID uint64) ([]models.Disposisi, error) {
	return s.disposisiRepo.GetByPimpinanID(pimpinanID)
}

func (s *PimpinanService) urgensiTinggi() ([]models.Pengaduan, error) {
	items, err := s.pengaduanRepo.GetAll()
	if err != nil {
		return nil, err
	}

	filtered := make([]models.Pengaduan, 0)
	for _, item := range items {
		if item.HasilAI != nil && strings.ToLower(item.HasilAI.Urgensi) == strings.ToLower("Tinggi") {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}
