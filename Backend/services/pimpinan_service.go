package services

import (
	"errors"
	"strconv"
	"strings"
	"time"

	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/repository"

	"gorm.io/gorm"
)

type PimpinanService struct {
	pengaduanRepo repository.PengaduanRepository
	disposisiRepo repository.DisposisiRepository
	unitRepo      repository.UnitRepository
}

func NewPimpinanService() *PimpinanService {
	return &PimpinanService{
		pengaduanRepo: repository.NewPengaduanRepository(),
		disposisiRepo: repository.NewDisposisiRepository(),
		unitRepo:      repository.NewUnitRepository(),
	}
}

func (s *PimpinanService) Dashboard() (*dto.DashboardPimpinanResponse, error) {
	all, err := s.pengaduanRepo.GetAll()
	if err != nil {
		return nil, err
	}
	resp := &dto.DashboardPimpinanResponse{TotalPengaduan: int64(len(all))}
	for _, item := range all {
		switch strings.ToLower(item.Status) {
		case strings.ToLower(StatusMenunggu), strings.ToLower(StatusMenungguDisposisi), strings.ToLower(StatusDiteruskanUnit):
			resp.BelumDikerjakan++
		case strings.ToLower(StatusDiproses):
			resp.SedangDiproses++
		case strings.ToLower(StatusSelesai):
			resp.Selesai++
		case strings.ToLower(StatusDitolak):
			resp.Ditolak++
		}
		if item.HasilAI != nil {
			switch strings.ToLower(item.HasilAI.Urgensi) {
			case "rendah":
				resp.UrgensiRendah++
			case "sedang":
				resp.UrgensiSedang++
			case "tinggi":
				resp.TotalUrgensiTinggi++
			}
		}
		if item.Disposisi == nil && item.HasilAI != nil && strings.EqualFold(item.HasilAI.Urgensi, "Tinggi") && strings.EqualFold(item.Status, StatusMenungguDisposisi) {
			resp.BelumDisposisi++
		} else if item.Disposisi != nil && item.HasilAI != nil && strings.EqualFold(item.HasilAI.Urgensi, "Tinggi") {
			resp.SudahDisposisi++
		}
	}
	return resp, nil
}

/*
The dashboard deliberately counts from the repository result instead of
frontend constants.  Keep the old urgent-disposition counters for existing
clients while exposing the complete monitoring summary above.
*/
func (s *PimpinanService) legacyDashboardCounters(items []models.Pengaduan, resp *dto.DashboardPimpinanResponse) {
	for _, item := range items {
		if item.Disposisi == nil {
			resp.BelumDisposisi++
		} else {
			resp.SudahDisposisi++
		}
	}
}

func (s *PimpinanService) GetUrgensiTinggi() ([]dto.PengaduanResponse, error) {
	items, err := s.urgensiTinggi()
	if err != nil {
		return nil, err
	}

	return mapPengaduanResponses(items), nil
}

type PimpinanHistoryFilter struct {
	Search     string
	Status     string
	KategoriID uint64
	UnitID     uint64
	From       *time.Time
	To         *time.Time
}

// GetRiwayatUrgensiTinggi returns every high-urgency complaint that actually
// reached the leadership queue, including records already disposed or closed.
func (s *PimpinanService) GetRiwayatUrgensiTinggi(filter PimpinanHistoryFilter) ([]dto.PengaduanResponse, error) {
	items, err := s.pengaduanRepo.GetAll()
	if err != nil {
		return nil, err
	}
	search := strings.ToLower(strings.TrimSpace(filter.Search))
	result := make([]dto.PengaduanResponse, 0)
	for index := range items {
		item := &items[index]
		if !s.isPimpinanHistoryItem(item) {
			continue
		}
		entryAt, _ := pimpinanEntryTime(item)
		if filter.From != nil && entryAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && !entryAt.Before(*filter.To) {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(item.Status, filter.Status) {
			continue
		}
		if filter.KategoriID != 0 && uint64(item.KategoriID) != filter.KategoriID {
			continue
		}
		if filter.UnitID != 0 && complaintUnitID(item) != filter.UnitID {
			continue
		}
		if search != "" && !complaintMatchesSearch(item, search) {
			continue
		}
		mapped := mapPengaduanResponse(item)
		mapped.TanggalMasukPimpinan = timePointer(entryAt)
		result = append(result, *mapped)
	}
	return result, nil
}

func (s *PimpinanService) GetRiwayatUrgensiTinggiDetail(id uint64) (*dto.PengaduanResponse, error) {
	item, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !s.isPimpinanHistoryItem(item) {
		return nil, errors.New("pengaduan tidak pernah masuk ke Pimpinan")
	}
	result := mapPengaduanResponse(item)
	if entryAt, ok := pimpinanEntryTime(item); ok {
		result.TanggalMasukPimpinan = timePointer(entryAt)
	}
	if err := attachCoordination(result, item.ID); err != nil {
		return nil, err
	}
	return result, nil
}

func attachCoordination(result *dto.PengaduanResponse, pengaduanID uint) error {
	var coordination []models.KoordinasiInternal
	if err := config.DB.Preload("Sender").Where("pengaduan_id = ?", pengaduanID).Order("created_at ASC").Find(&coordination).Error; err != nil {
		return err
	}
	for _, message := range coordination {
		mapped := mapCoordination(message)
		if message.Lampiran != "" {
			mapped.LampiranURL = "/api/pengaduan/" + formatUint(pengaduanID) + "/koordinasi/" + formatUint(uint(message.ID)) + "/lampiran"
		}
		result.KoordinasiInternal = append(result.KoordinasiInternal, mapped)
	}
	return nil
}

func (s *PimpinanService) isPimpinanHistoryItem(item *models.Pengaduan) bool {
	return item.HasilAI != nil &&
		strings.EqualFold(item.HasilAI.Urgensi, "Tinggi") &&
		item.Validasi != nil &&
		strings.EqualFold(item.Validasi.StatusValidasi, "Diterima") &&
		hasPimpinanEntry(item)
}

func hasPimpinanEntry(item *models.Pengaduan) bool {
	for _, history := range item.RiwayatStatus {
		if strings.EqualFold(history.StatusBaru, StatusMenungguDisposisi) {
			return true
		}
	}
	return false
}

func pimpinanEntryTime(item *models.Pengaduan) (time.Time, bool) {
	var result time.Time
	for _, history := range item.RiwayatStatus {
		if strings.EqualFold(history.StatusBaru, StatusMenungguDisposisi) && (result.IsZero() || history.CreatedAt.Before(result)) {
			result = history.CreatedAt
		}
	}
	return result, !result.IsZero()
}

func complaintUnitID(item *models.Pengaduan) uint64 {
	if item.UnitID != nil {
		return uint64(*item.UnitID)
	}
	if item.Disposisi != nil {
		return uint64(item.Disposisi.UnitID)
	}
	return 0
}

func complaintMatchesSearch(item *models.Pengaduan, search string) bool {
	student := ""
	if item.User.Mahasiswa != nil {
		student = item.User.Mahasiswa.NIM
	}
	value := strings.ToLower(strings.Join([]string{item.KodeTiket, item.Judul, item.User.NamaLengkap, student}, " "))
	return strings.Contains(value, search)
}

func timePointer(value time.Time) *time.Time { return &value }

func formatUint(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}

func (s *PimpinanService) GetPengaduan(id uint64) (*dto.PengaduanResponse, error) {
	item, err := s.pengaduanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	result := mapPengaduanResponse(item)
	if err := attachCoordination(result, item.ID); err != nil {
		return nil, err
	}
	if entryAt, ok := pimpinanEntryTime(item); ok {
		result.TanggalMasukPimpinan = timePointer(entryAt)
	}
	return result, nil
}

func (s *PimpinanService) Monitoring() ([]dto.PengaduanResponse, error) {
	items, err := s.pengaduanRepo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]dto.PengaduanResponse, 0)
	for index := range items {
		item := &items[index]
		if !isMonitoringItem(item) {
			continue
		}
		mapped := mapPengaduanResponse(item)
		if entryAt, ok := pimpinanEntryTime(item); ok {
			mapped.TanggalMasukPimpinan = timePointer(entryAt)
		}
		result = append(result, *mapped)
	}
	return result, nil
}

func isMonitoringItem(item *models.Pengaduan) bool {
	if item.Disposisi != nil || item.UnitID != nil {
		return true
	}
	switch strings.ToLower(item.Status) {
	case strings.ToLower(StatusMenungguDisposisi), strings.ToLower(StatusDiteruskanUnit), strings.ToLower(StatusDiproses), strings.ToLower(StatusSelesai):
		return true
	default:
		return item.HasilAI != nil &&
			strings.EqualFold(item.HasilAI.Urgensi, "Tinggi") &&
			item.Validasi != nil &&
			strings.EqualFold(item.Validasi.StatusValidasi, "Diterima")
	}
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
	unit, err := s.unitRepo.GetWithActiveKasubagByID(uint64(req.UnitID))
	if err != nil || unit == nil {
		return errors.New("unit tujuan tidak memiliki Kasubag aktif")
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

func (s *PimpinanService) GetAssignmentUnits() ([]models.Unit, error) {
	return s.unitRepo.GetWithActiveKasubag()
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
