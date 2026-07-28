package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func (s *PengaduanService) GetCategories() ([]models.KategoriPengaduan, error) {
	return s.kategoriRepo.GetAll()
}

func (s *PengaduanService) Create(userID uint, req dto.CreatePengaduanRequest) (*dto.PengaduanResponse, error) {
	if _, err := s.kategoriRepo.GetByID(uint64(req.KategoriID)); err != nil {
		return nil, errors.New("kategori pengaduan tidak ditemukan")
	}

	analysis, analysisErr := s.aiService.Analyze(dto.AIRequest{Deskripsi: strings.TrimSpace(req.Deskripsi)})
	if analysisErr != nil {
		log.Printf("analisis AI pengaduan tertunda: %v", analysisErr)
	}

	pengaduan := &models.Pengaduan{
		KodeTiket:        utils.GenerateTicketCode(),
		UserID:           userID,
		KategoriID:       req.KategoriID,
		Judul:            req.Judul,
		Deskripsi:        req.Deskripsi,
		Lampiran:         req.Lampiran,
		LampiranNamaAsli: optionalString(req.LampiranNamaAsli),
		LampiranMimeType: optionalString(req.LampiranMimeType),
		LampiranUkuran:   optionalInt64(req.LampiranUkuran),
		Status:           "Menunggu Verifikasi",
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pengaduan).Error; err != nil {
			return err
		}
		if analysis != nil {
			if err := tx.Create(hasilAIFromResponse(pengaduan.ID, analysis)).Error; err != nil {
				return err
			}
		}
		var admins []models.User
		if err := tx.Where("role = ? AND is_active = ?", utils.RoleAdminFakultas, true).Find(&admins).Error; err != nil {
			return err
		}
		for _, admin := range admins {
			if err := createNotification(tx, admin.ID, pengaduan.ID, "Pengaduan Baru", "Pengaduan "+pengaduan.KodeTiket+" menunggu verifikasi."); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if analysis != nil {
		pengaduan.HasilAI = hasilAIFromResponse(pengaduan.ID, analysis)
	}

	return mapPengaduanResponse(pengaduan), nil
}

func hasilAIFromResponse(pengaduanID uint, analysis *dto.AIResponse) *models.HasilAI {
	tokens, err := json.Marshal(analysis.Tokens)
	if err != nil || len(tokens) == 0 {
		tokens = []byte("[]")
	}
	positive, negative := scoreComponents(analysis.Score)
	return &models.HasilAI{
		PengaduanID:        pengaduanID,
		CleanedText:        analysis.CleanedText,
		Tokens:             models.JSONB(tokens),
		SkorPositif:        &positive,
		SkorNegatif:        &negative,
		SkorSentimen:       analysis.Score,
		Sentimen:           analysis.Sentimen,
		PenjelasanSentimen: fmt.Sprintf("Skor sentimen sebesar %d berada %s 0 sehingga dikategorikan sebagai %s.", analysis.Score, scoreRelation(analysis.Score), analysis.Sentimen),
		DetailSkor:         models.JSONB([]byte("[]")),
		Urgensi:            analysis.Urgensi,
		DasarUrgensi:       "",
	}
}

// The final database requires a non-null sign decomposition whose sum equals
// skor_sentimen. This preserves the actual aggregate score without pretending
// to know token weights that the AI response does not provide.
func scoreComponents(score int) (int, int) {
	if score >= 0 {
		return score, 0
	}
	return 0, score
}

func scoreRelation(score int) string {
	if score > 0 {
		return "di atas"
	}
	if score < 0 {
		return "di bawah"
	}
	return "sama dengan"
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
func int64Value(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
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
	if strings.EqualFold(role, utils.RoleKasubag) {
		if !userBelongsToComplaintUnit(userID, pengaduan.UnitID) {
			return nil, ErrForbidden
		}
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
	if strings.EqualFold(role, utils.RoleKasubag) {
		if !userBelongsToComplaintUnit(userID, pengaduan.UnitID) {
			return nil, ErrForbidden
		}
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
			"kategori_id":        pengaduan.KategoriID,
			"judul":              pengaduan.Judul,
			"deskripsi":          pengaduan.Deskripsi,
			"lampiran":           pengaduan.Lampiran,
			"lampiran_nama_asli": pengaduan.LampiranNamaAsli,
			"lampiran_mime_type": pengaduan.LampiranMimeType,
			"lampiran_ukuran":    pengaduan.LampiranUkuran,
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
			return tx.Create(hasilAIFromResponse(uint(id), analysis)).Error
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

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.ResponPengaduan{
			PengaduanID: uint(pengaduanID),
			UserID:      userID,
			Pesan:       pesan,
		}).Error; err != nil {
			return err
		}
		if !isMahasiswaRole(role) {
			return createNotification(tx, pengaduan.UserID, pengaduan.ID, "Balasan Baru", "Ada balasan baru untuk aduan "+pengaduan.KodeTiket)
		}
		return nil
	}); err != nil {
		return err
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

	if !strings.EqualFold(pengaduan.Status, StatusDiproses) {
		return errors.New("pengaduan hanya dapat diselesaikan saat Diproses")
	}
	now := time.Now()
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND user_id = ? AND status = ?", id, userID, StatusDiproses).Updates(map[string]interface{}{"status": StatusSelesai, "tanggal_selesai": &now}).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, pengaduan.ID, uint(userID), pengaduan.Status, StatusSelesai, "Diselesaikan oleh mahasiswa"); err != nil {
			return err
		}
		return createNotification(tx, pengaduan.UserID, pengaduan.ID, "Pengaduan Selesai", "Pengaduan "+pengaduan.KodeTiket+" telah diselesaikan.")
	})
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
		ID:               pengaduan.ID,
		KodeTiket:        pengaduan.KodeTiket,
		UserID:           pengaduan.UserID,
		KategoriID:       pengaduan.KategoriID,
		UnitID:           pengaduan.UnitID,
		Judul:            pengaduan.Judul,
		Deskripsi:        pengaduan.Deskripsi,
		Lampiran:         pengaduan.Lampiran,
		LampiranNamaAsli: stringValue(pengaduan.LampiranNamaAsli),
		LampiranMimeType: stringValue(pengaduan.LampiranMimeType),
		LampiranUkuran:   int64Value(pengaduan.LampiranUkuran),
		Status:           pengaduan.Status,
		AIStatus:         "pending",
		CreatedAt:        pengaduan.CreatedAt,
	}
	if pengaduan.Lampiran != "" {
		response.LampiranURL = fmt.Sprintf("/api/pengaduan/%d/lampiran", pengaduan.ID)
	}
	if pengaduan.Kategori.ID != 0 {
		response.Kategori = &dto.KategoriResponse{ID: pengaduan.Kategori.ID, Nama: pengaduan.Kategori.Nama}
	}
	if pengaduan.Unit.ID != 0 {
		response.Unit = &dto.UnitResponse{ID: pengaduan.Unit.ID, NamaUnit: pengaduan.Unit.NamaUnit}
	}

	if pengaduan.HasilAI != nil {
		skor := pengaduan.HasilAI.SkorSentimen
		response.SkorSentimen = &skor
		response.Sentimen = pengaduan.HasilAI.Sentimen
		response.Urgensi = pengaduan.HasilAI.Urgensi
		response.AIStatus = "success"
	}
	if pengaduan.Validasi != nil {
		response.Validasi = &dto.ValidasiResponse{ID: pengaduan.Validasi.ID, AdminFakultasID: pengaduan.Validasi.AdminFakultasID, StatusValidasi: pengaduan.Validasi.StatusValidasi, Catatan: pengaduan.Validasi.Catatan}
	}
	if pengaduan.Disposisi != nil {
		response.Disposisi = &dto.DisposisiResponse{ID: pengaduan.Disposisi.ID, PimpinanID: pengaduan.Disposisi.PimpinanID, UnitID: pengaduan.Disposisi.UnitID, Catatan: pengaduan.Disposisi.Catatan}
	}
	if len(pengaduan.RiwayatStatus) > 0 {
		response.RiwayatStatus = make([]dto.RiwayatStatusResponse, 0, len(pengaduan.RiwayatStatus))
		for _, history := range pengaduan.RiwayatStatus {
			response.RiwayatStatus = append(response.RiwayatStatus, dto.RiwayatStatusResponse{ID: history.ID, ChangedBy: history.ChangedBy, StatusLama: history.StatusLama, StatusBaru: history.StatusBaru, Catatan: history.Catatan, CreatedAt: history.CreatedAt})
		}
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
				ID:               item.ID,
				PengaduanID:      item.PengaduanID,
				UserID:           item.UserID,
				Pesan:            item.Pesan,
				Lampiran:         item.Lampiran,
				LampiranNamaAsli: stringValue(item.LampiranNamaAsli),
				LampiranMimeType: stringValue(item.LampiranMimeType),
				LampiranUkuran:   int64Value(item.LampiranUkuran),
				CreatedAt:        item.CreatedAt,
			}
			if item.Lampiran != "" {
				respon.LampiranURL = fmt.Sprintf("/api/pengaduan/%d/respon/%d/lampiran", pengaduan.ID, item.ID)
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

type ProtectedAttachment struct {
	Path        string
	ContentType string
}

func (s *PengaduanService) GetComplaintAttachment(userID uint64, role string, complaintID uint64) (*ProtectedAttachment, error) {
	complaint, err := s.repo.GetByID(complaintID)
	if err != nil {
		return nil, err
	}
	if err := s.assertAttachmentAccess(complaint, userID, role); err != nil {
		return nil, err
	}
	return protectedAttachment(complaint.Lampiran)
}

func (s *PengaduanService) GetResponseAttachment(userID uint64, role string, complaintID, responseID uint64) (*ProtectedAttachment, error) {
	complaint, err := s.repo.GetByID(complaintID)
	if err != nil {
		return nil, err
	}
	if err := s.assertAttachmentAccess(complaint, userID, role); err != nil {
		return nil, err
	}
	response, err := s.responRepo.GetByID(responseID)
	if err != nil {
		return nil, err
	}
	if uint64(response.PengaduanID) != complaintID {
		return nil, gorm.ErrRecordNotFound
	}
	return protectedAttachment(response.Lampiran)
}

func (s *PengaduanService) assertAttachmentAccess(complaint *models.Pengaduan, userID uint64, role string) error {
	switch utils.CanonicalRole(role) {
	case utils.RoleMahasiswa:
		if uint64(complaint.UserID) != userID {
			return ErrForbidden
		}
	case utils.RoleKasubag:
		if !userBelongsToComplaintUnit(userID, complaint.UnitID) {
			return ErrForbidden
		}
	case utils.RoleAdminFakultas, utils.RolePimpinan:
		// These roles are scoped to the faculty/application as in the existing API.
	default:
		return ErrForbidden
	}
	return nil
}

func protectedAttachment(filename string) (*ProtectedAttachment, error) {
	path, err := utils.ResolveUploadedFile(filename)
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	header := make([]byte, 512)
	read, err := file.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return &ProtectedAttachment{Path: path, ContentType: http.DetectContentType(header[:read])}, nil
}

func isMahasiswaRole(role string) bool {
	return strings.EqualFold(role, "mahasiswa")
}

func userBelongsToComplaintUnit(userID uint64, complaintUnitID *uint) bool {
	var user models.User
	if err := config.DB.Select("unit_id").First(&user, userID).Error; err != nil || user.UnitID == nil || complaintUnitID == nil {
		return false
	}
	return *user.UnitID == *complaintUnitID
}
