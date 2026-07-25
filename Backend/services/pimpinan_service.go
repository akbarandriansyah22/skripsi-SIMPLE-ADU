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

func (s *PimpinanService) GetPengaduan(id uint64) (*dto.PengaduanResponse, error) {
	item, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return mapPengaduanResponse(item), nil
}

func (s *PimpinanService) Monitoring() ([]dto.PengaduanResponse, error) {
	items, err := s.pengaduanRepo.GetAll()
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
	if !strings.EqualFold(pengaduan.Status, StatusMenungguDisposisi) {
		return errors.New("pengaduan belum diteruskan oleh admin")
	}
	if strings.TrimSpace(req.Catatan) == "" {
		return errors.New("catatan disposisi wajib diisi")
	}
	unit, err := repository.NewUnitRepository().GetByID(uint64(req.UnitID))
	if err != nil || unit == nil {
		return errors.New("unit disposisi tidak ditemukan")
	}
	if unit.NamaUnit != "Akademik" && unit.NamaUnit != "Sarana dan Prasarana" {
		return errors.New("disposisi hanya dapat diarahkan ke unit Akademik atau Sarana dan Prasarana")
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
		unitID := req.UnitID
		if err := tx.Create(&models.Disposisi{PengaduanID: pengaduan.ID, PimpinanID: pimpinanID, UnitID: unitID, Catatan: strings.TrimSpace(req.Catatan)}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND status = ?", pengaduanID, StatusMenungguDisposisi).Updates(map[string]interface{}{"unit_id": unitID, "status": StatusDiteruskanUnit}).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, pengaduan.ID, pimpinanID, pengaduan.Status, StatusDiteruskanUnit, req.Catatan); err != nil {
			return err
		}
		var kasubags []models.User
		if err := tx.Where("role = ? AND unit_id = ? AND is_active = ?", "kasubag", unitID, true).Find(&kasubags).Error; err != nil {
			return err
		}
		for _, kasubag := range kasubags {
			if err := createNotification(tx, kasubag.ID, pengaduan.ID, "Aduan Diteruskan ke Unit", "Aduan "+pengaduan.KodeTiket+" telah didisposisikan ke unit Anda."); err != nil {
				return err
			}
		}
		return createNotification(tx, pengaduan.UserID, pengaduan.ID, "Disposisi Pengaduan", "Pengaduan "+pengaduan.KodeTiket+" telah mendapat disposisi pimpinan.")
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
		if item.HasilAI != nil && strings.EqualFold(item.HasilAI.Urgensi, "Tinggi") && strings.EqualFold(item.Status, StatusDiteruskan) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}
