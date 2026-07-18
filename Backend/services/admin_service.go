package services

import (
	"errors"
	"log"
	"strings"

	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/utils"

	"gorm.io/gorm"
)

type AdminService struct {
	pengaduanRepo repository.PengaduanRepository
	hasilAIRepo   repository.HasilAIRepository
	unitRepo      repository.UnitRepository
	disposisiRepo repository.DisposisiRepository
	aiService     *AIService
}

func NewAdminService() *AdminService {
	return &AdminService{
		pengaduanRepo: repository.NewPengaduanRepository(),
		hasilAIRepo:   repository.NewHasilAIRepository(),
		unitRepo:      repository.NewUnitRepository(),
		disposisiRepo: repository.NewDisposisiRepository(),
		aiService:     NewAIService(),
	}
}

func (s *AdminService) Dashboard() (*dto.DashboardAdminResponse, error) {
	items, err := s.pengaduanRepo.GetAll()
	if err != nil {
		return nil, err
	}

	resp := &dto.DashboardAdminResponse{
		TotalPengaduan: int64(len(items)),
	}

	for _, item := range items {
		switch strings.ToLower(item.Status) {
		case strings.ToLower("Menunggu Verifikasi"):
			resp.Menunggu++
		case strings.ToLower("Diproses"):
			resp.Diproses++
		case strings.ToLower("Selesai"):
			resp.Selesai++
		case strings.ToLower("Ditolak"):
			resp.Ditolak++
		}
	}

	return resp, nil
}

func (s *AdminService) GetAllPengaduan() ([]dto.PengaduanResponse, error) {
	pengaduanService := NewPengaduanService()
	return pengaduanService.GetAll()
}

func (s *AdminService) GetPengaduanByID(id uint64) (*dto.PengaduanResponse, error) {
	pengaduanService := NewPengaduanService()
	return pengaduanService.GetByID(id)
}

func (s *AdminService) UpdateStatus(id uint64, req dto.UpdateStatusRequest) error {
	if req.Status == "" {
		return errors.New("status wajib diisi")
	}

	if !isAllowedPengaduanStatus(req.Status) {
		return errors.New("status tidak valid")
	}

	pengaduan, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return err
	}

	if !isValidStatusTransition(pengaduan.Status, req.Status) {
		return errors.New("transisi status tidak diizinkan")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
			return err
		}
		return tx.Create(&models.Notifikasi{UserID: pengaduan.UserID, Judul: "Status Aduan Diperbarui", Isi: "Status aduan " + pengaduan.KodeTiket + " berubah menjadi " + req.Status}).Error
	})
}

func (s *AdminService) AssignUnit(id uint64, req dto.AssignUnitRequest) error {
	if _, err := s.unitRepo.GetByID(uint64(req.UnitID)); err != nil {
		return errors.New("unit tidak ditemukan")
	}

	pengaduan, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return err
	}

	if !strings.EqualFold(pengaduan.Status, StatusMenunggu) {
		return errors.New("unit hanya dapat ditetapkan saat menunggu verifikasi")
	}
	if pengaduan.HasilAI != nil && strings.EqualFold(pengaduan.HasilAI.Urgensi, "Tinggi") {
		return errors.New("pengaduan urgensi Tinggi harus diteruskan ke pimpinan")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ?", id).Updates(map[string]interface{}{"unit_id": req.UnitID, "status": StatusDiproses}).Error; err != nil {
			return err
		}
		return tx.Create(&models.Notifikasi{UserID: pengaduan.UserID, Judul: "Pengaduan Mulai Diproses", Isi: "Pengaduan " + pengaduan.KodeTiket + " telah diteruskan ke unit terkait."}).Error
	})
}

func (s *AdminService) ForwardToPimpinan(id uint64) error {
	pengaduan, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return err
	}

	if pengaduan.HasilAI == nil || strings.ToLower(pengaduan.HasilAI.Urgensi) != strings.ToLower("Tinggi") {
		return errors.New("hanya pengaduan urgensi Tinggi yang diteruskan ke pimpinan")
	}
	if !strings.EqualFold(pengaduan.Status, StatusMenunggu) {
		return errors.New("pengaduan hanya dapat diteruskan setelah verifikasi admin")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ?", id).Update("status", StatusDiteruskan).Error; err != nil {
			return err
		}
		var pimpinan []models.User
		if err := tx.Where("role = ?", utils.RolePimpinan).Find(&pimpinan).Error; err != nil {
			return err
		}
		for _, user := range pimpinan {
			if err := tx.Create(&models.Notifikasi{UserID: user.ID, Judul: "Pengaduan Urgensi Tinggi", Isi: "Pengaduan " + pengaduan.KodeTiket + " menunggu disposisi."}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *AdminService) ReanalyzeAI(id uint64) (*dto.ReanalyzeAIResponse, error) {
	pengaduan, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	analysis, err := s.aiService.Analyze(dto.AIRequest{Deskripsi: pengaduan.Deskripsi})
	if err != nil {
		log.Printf("retry analisis AI pengaduan %d gagal: %v", id, err)
		return nil, err
	}

	hasil := &models.HasilAI{
		PengaduanID:  pengaduan.ID,
		SkorSentimen: analysis.Score,
		Sentimen:     analysis.Sentimen,
		Urgensi:      analysis.Urgensi,
	}
	if err := s.hasilAIRepo.UpsertByPengaduanID(hasil); err != nil {
		return nil, err
	}

	return &dto.ReanalyzeAIResponse{
		PengaduanID:  pengaduan.ID,
		SkorSentimen: analysis.Score,
		Sentimen:     analysis.Sentimen,
		Urgensi:      analysis.Urgensi,
		AIStatus:     "success",
	}, nil
}

func (s *AdminService) GetUnits() ([]models.Unit, error) {
	return s.unitRepo.GetAll()
}

func isAllowedPengaduanStatus(status string) bool {
	switch strings.ToLower(status) {
	case strings.ToLower(StatusMenunggu), strings.ToLower(StatusDiproses), strings.ToLower(StatusSelesai), strings.ToLower(StatusDitolak), strings.ToLower(StatusDiteruskan):
		return true
	default:
		return false
	}
}
