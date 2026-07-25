package services

import (
	"errors"
	"log"
	"strings"
	"time"

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
	return errors.New("Admin Fakultas mengubah status melalui validasi atau penetapan unit")
}

func (s *AdminService) AssignUnit(adminID, id uint64, req dto.AssignUnitRequest) error {
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
	if pengaduan.Validasi == nil || !strings.EqualFold(pengaduan.Validasi.StatusValidasi, "Diterima") {
		return errors.New("pengaduan harus diterima melalui validasi admin terlebih dahulu")
	}
	if pengaduan.HasilAI != nil && strings.EqualFold(pengaduan.HasilAI.Urgensi, "Tinggi") {
		return errors.New("pengaduan urgensi Tinggi harus diteruskan ke pimpinan")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND status = ?", id, StatusMenunggu).Updates(map[string]interface{}{"unit_id": req.UnitID, "status": StatusDiteruskanUnit}).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, pengaduan.ID, uint(adminID), pengaduan.Status, StatusDiteruskanUnit, "Diteruskan oleh Admin Fakultas"); err != nil {
			return err
		}
		var kasubags []models.User
		if err := tx.Where("role = ? AND unit_id = ? AND is_active = ?", utils.RoleKasubag, req.UnitID, true).Find(&kasubags).Error; err != nil {
			return err
		}
		for _, kasubag := range kasubags {
			if err := createNotification(tx, kasubag.ID, pengaduan.ID, "Aduan Diteruskan ke Unit", "Aduan "+pengaduan.KodeTiket+" telah diteruskan ke unit Anda."); err != nil {
				return err
			}
		}
		return createNotification(tx, pengaduan.UserID, pengaduan.ID, "Aduan Diteruskan ke Unit", "Aduan "+pengaduan.KodeTiket+" telah diteruskan ke unit terkait.")
	})
}

func (s *AdminService) ForwardToPimpinan(adminID, id uint64) error {
	pengaduan, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return err
	}

	if pengaduan.HasilAI == nil || !strings.EqualFold(pengaduan.HasilAI.Urgensi, "Tinggi") {
		return errors.New("hanya pengaduan urgensi Tinggi yang diteruskan ke pimpinan")
	}
	if !strings.EqualFold(pengaduan.Status, StatusMenunggu) && !strings.EqualFold(pengaduan.Status, StatusMenungguDisposisi) {
		return errors.New("pengaduan hanya dapat diteruskan setelah verifikasi admin")
	}
	if pengaduan.Validasi == nil || !strings.EqualFold(pengaduan.Validasi.StatusValidasi, "Diterima") {
		return errors.New("pengaduan harus diterima melalui validasi admin terlebih dahulu")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if strings.EqualFold(pengaduan.Status, StatusMenungguDisposisi) {
			var pimpinan []models.User
			if err := tx.Where("role = ? AND is_active = ?", utils.RolePimpinan, true).Find(&pimpinan).Error; err != nil {
				return err
			}
			for _, user := range pimpinan {
				if err := createNotification(tx, user.ID, pengaduan.ID, "Pengaduan Urgensi Tinggi", "Pengaduan "+pengaduan.KodeTiket+" menunggu disposisi."); err != nil {
					return err
				}
			}
			return nil
		}
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND status = ?", id, StatusMenunggu).Update("status", StatusMenungguDisposisi).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, pengaduan.ID, uint(adminID), pengaduan.Status, StatusMenungguDisposisi, "Menunggu disposisi pimpinan"); err != nil {
			return err
		}
		var pimpinan []models.User
		if err := tx.Where("role = ?", utils.RolePimpinan).Find(&pimpinan).Error; err != nil {
			return err
		}
		for _, user := range pimpinan {
			if err := createNotification(tx, user.ID, pengaduan.ID, "Pengaduan Urgensi Tinggi", "Pengaduan "+pengaduan.KodeTiket+" menunggu disposisi."); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *AdminService) ValidatePengaduan(adminID, id uint64, req dto.ValidatePengaduanRequest) error {
	status := strings.TrimSpace(req.StatusValidasi)
	if !strings.EqualFold(status, "Diterima") && !strings.EqualFold(status, "Ditolak") {
		return errors.New("status validasi harus Diterima atau Ditolak")
	}
	pengaduan, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return err
	}
	if !strings.EqualFold(pengaduan.Status, StatusMenunggu) && !strings.EqualFold(pengaduan.Status, StatusMenungguDisposisi) {
		return errors.New("pengaduan tidak berada pada tahap validasi")
	}
	nextStatus := StatusMenunggu
	if strings.EqualFold(status, "Ditolak") {
		nextStatus = StatusDitolak
	} else if pengaduan.HasilAI != nil && strings.EqualFold(pengaduan.HasilAI.Urgensi, "Tinggi") {
		nextStatus = StatusMenungguDisposisi
	}
	var admin models.User
	if err := config.DB.Select("id", "role").First(&admin, adminID).Error; err != nil || utils.CanonicalRole(admin.Role) != utils.RoleAdminFakultas {
		return errors.New("validasi hanya dapat dilakukan oleh Admin Fakultas")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var validation models.ValidasiPengaduan
		err := tx.Where("pengaduan_id = ?", id).First(&validation).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			validation = models.ValidasiPengaduan{PengaduanID: uint(id)}
		} else if err != nil {
			return err
		}
		validation.AdminFakultasID = uint(adminID)
		validation.StatusValidasi = status
		validation.Catatan = strings.TrimSpace(req.Catatan)
		validation.ValidatedAt = time.Now()
		if validation.ID == 0 {
			if err := tx.Create(&validation).Error; err != nil {
				return err
			}
		} else if err := tx.Save(&validation).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Pengaduan{}).Where("id = ?", id).Update("status", nextStatus).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, pengaduan.ID, uint(adminID), pengaduan.Status, nextStatus, req.Catatan); err != nil {
			return err
		}
		if nextStatus == StatusMenungguDisposisi {
			var pimpinan []models.User
			if err := tx.Where("role = ? AND is_active = ?", utils.RolePimpinan, true).Find(&pimpinan).Error; err != nil {
				return err
			}
			for _, user := range pimpinan {
				if err := createNotification(tx, user.ID, pengaduan.ID, "Pengaduan Urgensi Tinggi", "Pengaduan "+pengaduan.KodeTiket+" menunggu disposisi."); err != nil {
					return err
				}
			}
		}
		return createNotification(tx, pengaduan.UserID, pengaduan.ID, "Validasi Aduan", "Aduan "+pengaduan.KodeTiket+" telah "+status+".")
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

	hasil := hasilAIFromResponse(pengaduan.ID, analysis)
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
	case strings.ToLower(StatusMenunggu), strings.ToLower(StatusMenungguDisposisi), strings.ToLower(StatusDiteruskanUnit), strings.ToLower(StatusDiproses), strings.ToLower(StatusSelesai), strings.ToLower(StatusDitolak):
		return true
	default:
		return false
	}
}
