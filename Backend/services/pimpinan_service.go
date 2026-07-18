package services

import (
	"errors"
	"strings"

	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/repository"

	"gorm.io/gorm"
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

	if pengaduan.HasilAI == nil || !strings.EqualFold(pengaduan.HasilAI.Urgensi, "Tinggi") {
		return errors.New("disposisi hanya dapat diberikan untuk urgensi Tinggi")
	}
	if !strings.EqualFold(pengaduan.Status, StatusDiteruskan) {
		return errors.New("pengaduan belum diteruskan oleh admin")
	}
	if _, err := s.disposisiRepo.GetByPengaduanID(pengaduanID); err == nil {
		return errors.New("disposisi untuk pengaduan ini sudah ada")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		var existing models.Disposisi
		if err := tx.Where("pengaduan_id = ?", pengaduanID).First(&existing).Error; err == nil {
			return errors.New("disposisi untuk pengaduan ini sudah ada")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Create(&models.Disposisi{PengaduanID: pengaduan.ID, PimpinanID: pimpinanID, Catatan: strings.TrimSpace(req.Catatan)}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Pengaduan{}).Where("id = ?", pengaduanID).Update("status", StatusDiproses).Error; err != nil {
			return err
		}
		return tx.Create(&models.Notifikasi{UserID: pengaduan.UserID, Judul: "Disposisi Pengaduan", Isi: "Pengaduan " + pengaduan.KodeTiket + " telah mendapat disposisi pimpinan."}).Error
	})
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
		if item.HasilAI != nil && strings.ToLower(item.HasilAI.Urgensi) == strings.ToLower("Tinggi") && strings.EqualFold(item.Status, StatusDiteruskan) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}
