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

var ErrForbidden = errors.New("akses ditolak")

type PengaduanService struct {
	repo         repository.PengaduanRepository
	kategoriRepo repository.KategoriRepository
	hasilAIRepo  repository.HasilAIRepository
	responRepo   repository.ResponRepository
	aiService    *AIService
}

func NewPengaduanService() *PengaduanService {
	return &PengaduanService{
		repo:         repository.NewPengaduanRepository(),
		kategoriRepo: repository.NewKategoriRepository(),
		hasilAIRepo:  repository.NewHasilAIRepository(),
		responRepo:   repository.NewResponRepository(),
		aiService:    NewAIService(),
	}
}

func (s *PengaduanService) Create(userID uint, req dto.CreatePengaduanRequest) (*dto.PengaduanResponse, error) {
	if _, err := s.kategoriRepo.GetByID(uint64(req.KategoriID)); err != nil {
		return nil, errors.New("kategori pengaduan tidak ditemukan")
	}

	pengaduan := &models.Pengaduan{
		KodeTiket:  utils.GenerateTicketCode(),
		UserID:     userID,
		KategoriID: req.KategoriID,
		Judul:      req.Judul,
		Deskripsi:  req.Deskripsi,
		Lampiran:   req.Lampiran,
		Status:     "Menunggu Verifikasi",
	}

	if err := config.DB.Create(pengaduan).Error; err != nil {
		return nil, err
	}

	if err := s.analyzeAndAttach(pengaduan); err != nil {
		log.Printf("analisis AI pengaduan %d tertunda: %v", pengaduan.ID, err)
	}

	return mapPengaduanResponse(pengaduan), nil
}

func (s *PengaduanService) analyzeAndAttach(pengaduan *models.Pengaduan) error {
	analysis, err := s.aiService.Analyze(dto.AIRequest{Deskripsi: pengaduan.Deskripsi})
	if err != nil {
		return err
	}

	hasilAI := &models.HasilAI{
		PengaduanID:  pengaduan.ID,
		Sentimen:     analysis.Sentimen,
		SkorSentimen: analysis.Score,
		Urgensi:      analysis.Urgensi,
	}

	if err := s.hasilAIRepo.UpsertByPengaduanID(hasilAI); err != nil {
		return err
	}

	pengaduan.HasilAI = hasilAI
	return nil
}

func (s *PengaduanService) GetByID(id uint64) (*dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return mapPengaduanResponse(pengaduan), nil
}

func (s *PengaduanService) GetByIDForRole(id uint64, userID uint64, role string) (*dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if isMahasiswaRole(role) && uint64(pengaduan.UserID) != userID {
		return nil, ErrForbidden
	}

	return mapPengaduanResponse(pengaduan), nil
}

func (s *PengaduanService) GetByKodeTiket(kode string) (*dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetByKodeTiket(kode)
	if err != nil {
		return nil, err
	}

	return mapPengaduanResponse(pengaduan), nil
}

func (s *PengaduanService) GetByKodeTiketForRole(kode string, userID uint64, role string) (*dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetByKodeTiket(kode)
	if err != nil {
		return nil, err
	}

	if isMahasiswaRole(role) && uint64(pengaduan.UserID) != userID {
		return nil, ErrForbidden
	}

	return mapPengaduanResponse(pengaduan), nil
}

func (s *PengaduanService) GetByUserID(userID uint64) ([]dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return mapPengaduanResponses(pengaduan), nil
}

func (s *PengaduanService) GetAll() ([]dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return mapPengaduanResponses(pengaduan), nil
}

func (s *PengaduanService) Update(userID uint64, id uint64, req dto.UpdatePengaduanRequest) (*dto.PengaduanResponse, error) {
	pengaduan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if uint64(pengaduan.UserID) != userID {
		return nil, errors.New("pengaduan bukan milik user")
	}

	if pengaduan.Status != "Menunggu Verifikasi" {
		return nil, errors.New("pengaduan hanya dapat diubah saat Menunggu Verifikasi")
	}
	originalDescription := pengaduan.Deskripsi

	newDescription := pengaduan.Deskripsi
	if req.KategoriID != 0 {
		if _, err := s.kategoriRepo.GetByID(uint64(req.KategoriID)); err != nil {
			return nil, errors.New("kategori pengaduan tidak ditemukan")
		}
		pengaduan.KategoriID = req.KategoriID
	}
	if req.Judul != "" {
		pengaduan.Judul = req.Judul
	}
	if req.Deskripsi != "" {
		newDescription = strings.TrimSpace(req.Deskripsi)
		if newDescription == "" {
			return nil, errors.New("deskripsi tidak boleh kosong")
		}
		pengaduan.Deskripsi = newDescription
	}
	if req.Lampiran != "" {
		pengaduan.Lampiran = req.Lampiran
	}

	descriptionChanged := req.Deskripsi != "" && strings.TrimSpace(req.Deskripsi) != originalDescription
	// Analyze the new text before committing it. If AI is unavailable, delete
	// the old result so the API exposes a pending state instead of stale data.
	var analysis *dto.AIResponse
	if descriptionChanged {
		var err error
		analysis, err = s.aiService.Analyze(dto.AIRequest{Deskripsi: strings.TrimSpace(req.Deskripsi)})
		if err != nil {
			log.Printf("analisis AI pengaduan %d tertunda setelah edit: %v", id, err)
		}
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ?", id).Updates(map[string]interface{}{
			"kategori_id": pengaduan.KategoriID,
			"judul":       pengaduan.Judul,
			"deskripsi":   pengaduan.Deskripsi,
			"lampiran":    pengaduan.Lampiran,
		}).Error; err != nil {
			return err
		}
		if !descriptionChanged {
			return nil
		}
		if err := tx.Where("pengaduan_id = ?", id).Delete(&models.HasilAI{}).Error; err != nil {
			return err
		}
		if analysis != nil {
			return tx.Create(&models.HasilAI{PengaduanID: uint(id), Sentimen: analysis.Sentimen, SkorSentimen: analysis.Score, Urgensi: analysis.Urgensi}).Error
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *PengaduanService) AddRespon(userID uint, role string, pengaduanID uint64, pesan string) error {
	if pesan == "" {
		return errors.New("pesan wajib diisi")
	}

	pengaduan, err := s.repo.GetByID(pengaduanID)
	if err != nil {
		return err
	}

	if isMahasiswaRole(role) && uint64(pengaduan.UserID) != uint64(userID) {
		return ErrForbidden
	}

	if err := s.responRepo.Create(&models.ResponPengaduan{
		PengaduanID: uint(pengaduanID),
		UserID:      userID,
		Pesan:       pesan,
	}); err != nil {
		return err
	}

	if !isMahasiswaRole(role) {
		_ = NewNotifikasiService().Create(
			pengaduan.UserID,
			"Balasan Baru dari Admin",
			"Ada balasan baru untuk aduan "+pengaduan.KodeTiket,
		)
	}

	return nil
}

func (s *PengaduanService) Finish(userID uint64, id uint64) error {
	pengaduan, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if uint64(pengaduan.UserID) != userID {
		return ErrForbidden
	}

	now := time.Now()
	pengaduan.Status = "Selesai"
	pengaduan.TanggalSelesai = &now

	return s.repo.Update(pengaduan)
}

func mapPengaduanResponses(items []models.Pengaduan) []dto.PengaduanResponse {
	responses := make([]dto.PengaduanResponse, 0, len(items))
	for i := range items {
		responses = append(responses, *mapPengaduanResponse(&items[i]))
	}

	return responses
}

func mapPengaduanResponse(pengaduan *models.Pengaduan) *dto.PengaduanResponse {
	response := &dto.PengaduanResponse{
		ID:         pengaduan.ID,
		KodeTiket:  pengaduan.KodeTiket,
		UserID:     pengaduan.UserID,
		KategoriID: pengaduan.KategoriID,
		UnitID:     pengaduan.UnitID,
		Judul:      pengaduan.Judul,
		Deskripsi:  pengaduan.Deskripsi,
		Lampiran:   pengaduan.Lampiran,
		Status:     pengaduan.Status,
		AIStatus:   "pending",
		CreatedAt:  pengaduan.CreatedAt,
	}

	if pengaduan.HasilAI != nil {
		skor := pengaduan.HasilAI.SkorSentimen
		response.SkorSentimen = &skor
		response.Sentimen = pengaduan.HasilAI.Sentimen
		response.Urgensi = pengaduan.HasilAI.Urgensi
		response.AIStatus = "success"
	}

	if pengaduan.User.ID != 0 {
		response.User = &dto.UserResponse{
			ID:          pengaduan.User.ID,
			NamaLengkap: pengaduan.User.NamaLengkap,
			Email:       pengaduan.User.Email,
			Role:        pengaduan.User.Role,
			IsActive:    pengaduan.User.IsActive,
		}
	}

	if len(pengaduan.ResponPengaduan) > 0 {
		response.ResponPengaduan = make([]dto.ResponPengaduanResponse, 0, len(pengaduan.ResponPengaduan))
		for _, item := range pengaduan.ResponPengaduan {
			respon := dto.ResponPengaduanResponse{
				ID:          item.ID,
				PengaduanID: item.PengaduanID,
				UserID:      item.UserID,
				Pesan:       item.Pesan,
				CreatedAt:   item.CreatedAt,
			}
			if item.User.ID != 0 {
				respon.User = &dto.UserResponse{
					ID:          item.User.ID,
					NamaLengkap: item.User.NamaLengkap,
					Email:       item.User.Email,
					Role:        item.User.Role,
					IsActive:    item.User.IsActive,
				}
			}
			response.ResponPengaduan = append(response.ResponPengaduan, respon)
		}
	}

	return response
}

func isMahasiswaRole(role string) bool {
	return strings.EqualFold(role, "mahasiswa")
}
