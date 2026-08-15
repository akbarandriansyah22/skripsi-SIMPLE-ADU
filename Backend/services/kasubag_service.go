package services

import (
	"errors"
	"strings"
	"time"

	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

type KasubagService struct {
	pengaduan *PengaduanService
}

func NewKasubagService() *KasubagService {
	return &KasubagService{pengaduan: NewPengaduanService()}
}

func (s *KasubagService) Dashboard(userID uint64) (map[string]interface{}, error) {
	items, unitID, err := s.unitComplaints(userID)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{"unit_id": unitID, "total": len(items), "diteruskan": 0, "diproses": 0, "selesai": 0}
	for _, item := range items {
		switch item.Status {
		case StatusDiteruskanUnit:
			result["diteruskan"] = result["diteruskan"].(int) + 1
		case StatusDiproses:
			result["diproses"] = result["diproses"].(int) + 1
		case StatusSelesai:
			result["selesai"] = result["selesai"].(int) + 1
		}
	}
	return result, nil
}

func (s *KasubagService) GetComplaints(userID uint64) ([]dto.PengaduanResponse, error) {
	items, _, err := s.unitComplaints(userID)
	if err != nil {
		return nil, err
	}
	return mapPengaduanResponses(items), nil
}

func (s *KasubagService) GetComplaint(userID, complaintID uint64) (*dto.PengaduanResponse, error) {
	complaint, err := s.pengaduan.GetByID(complaintID)
	if err != nil {
		return nil, err
	}
	if err := s.assertUnit(userID, complaint.UnitID); err != nil {
		return nil, err
	}
	return complaint, nil
}

func (s *KasubagService) StartProcess(userID, complaintID uint64) error {
	complaint, err := s.pengaduan.GetByID(complaintID)
	if err != nil {
		return err
	}
	if err := s.assertUnit(userID, complaint.UnitID); err != nil {
		return err
	}
	if !strings.EqualFold(complaint.Status, StatusDiteruskanUnit) {
		return errors.New("aduan belum diteruskan ke unit atau sudah diproses")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND unit_id = ? AND status = ?", complaintID, *complaint.UnitID, StatusDiteruskanUnit).Update("status", StatusDiproses).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, complaint.ID, uint(userID), complaint.Status, StatusDiproses, "Mulai diproses oleh Kasubag"); err != nil {
			return err
		}
		return createNotification(tx, complaint.UserID, complaint.ID, "Aduan Diproses", "Aduan "+complaint.KodeTiket+" mulai diproses oleh unit terkait.")
	})
}

func (s *KasubagService) AddResponse(userID, complaintID uint64, message, attachment, originalName, mimeType string, attachmentSize int64) error {
	complaint, err := s.pengaduan.GetByID(complaintID)
	if err != nil {
		return err
	}
	if err := s.assertUnit(userID, complaint.UnitID); err != nil {
		return err
	}
	if !strings.EqualFold(complaint.Status, StatusDiproses) {
		return errors.New("tindak lanjut hanya dapat diberikan saat aduan Diproses")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.ResponPengaduan{PengaduanID: uint(complaintID), UserID: uint(userID), Pesan: strings.TrimSpace(message), Lampiran: attachment, LampiranNamaAsli: optionalString(originalName), LampiranMimeType: optionalString(mimeType), LampiranUkuran: optionalInt64(attachmentSize)}).Error; err != nil {
			return err
		}
		return createNotification(tx, complaint.UserID, complaint.ID, "Tindak Lanjut Baru", "Ada tindak lanjut baru untuk aduan "+complaint.KodeTiket+".")
	})
}

func (s *KasubagService) Finish(userID, complaintID uint64) error {
	complaint, err := s.pengaduan.GetByID(complaintID)
	if err != nil {
		return err
	}
	if err := s.assertUnit(userID, complaint.UnitID); err != nil {
		return err
	}
	if !strings.EqualFold(complaint.Status, StatusDiproses) {
		return errors.New("aduan harus berstatus Diproses sebelum diselesaikan")
	}
	now := time.Now()
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND unit_id = ? AND status = ?", complaintID, *complaint.UnitID, StatusDiproses).Updates(map[string]interface{}{"status": StatusSelesai, "tanggal_selesai": &now}).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, complaint.ID, uint(userID), complaint.Status, StatusSelesai, "Diselesaikan oleh Kasubag"); err != nil {
			return err
		}
		return createNotification(tx, complaint.UserID, complaint.ID, "Pengaduan Selesai", "Aduan "+complaint.KodeTiket+" telah diselesaikan.")
	})
}

func (s *KasubagService) ReturnToAdmin(userID, complaintID uint64, reason string) error {
	complaint, err := s.pengaduan.GetByID(complaintID)
	if err != nil {
		return err
	}
	if err := s.assertUnit(userID, complaint.UnitID); err != nil {
		return err
	}
	if !strings.EqualFold(complaint.Status, StatusDiteruskanUnit) && !strings.EqualFold(complaint.Status, StatusDiproses) {
		return errors.New("aduan hanya dapat dikembalikan saat Diteruskan ke Unit atau Diproses")
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("alasan pengembalian wajib diisi")
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Pengaduan{}).Where("id = ? AND unit_id = ? AND status IN (?, ?)", complaintID, *complaint.UnitID, StatusDiteruskanUnit, StatusDiproses).Updates(map[string]interface{}{"unit_id": nil, "status": StatusMenunggu}).Error; err != nil {
			return err
		}
		if err := tx.Where("pengaduan_id = ?", complaint.ID).Delete(&models.Disposisi{}).Error; err != nil {
			return err
		}
		if err := recordStatusChange(tx, complaint.ID, uint(userID), complaint.Status, StatusMenunggu, reason); err != nil {
			return err
		}
		var admins []models.User
		if err := tx.Where("role = ? AND is_active = ?", utils.RoleAdminFakultas, true).Find(&admins).Error; err != nil {
			return err
		}
		for _, admin := range admins {
			if err := createNotification(tx, admin.ID, complaint.ID, "Aduan Dikembalikan", "Aduan "+complaint.KodeTiket+" dikembalikan oleh unit: "+strings.TrimSpace(reason)); err != nil {
				return err
			}
		}
		return createNotification(tx, complaint.UserID, complaint.ID, "Aduan Dikembalikan ke Admin", "Aduan "+complaint.KodeTiket+" dikembalikan ke Admin Fakultas.")
	})
}

func (s *KasubagService) unitComplaints(userID uint64) ([]models.Pengaduan, uint, error) {
	unitID, err := s.userUnitID(userID)
	if err != nil {
		return nil, 0, err
	}
	var items []models.Pengaduan
	err = config.DB.Where("unit_id = ?", unitID).
		Preload("User").Preload("Kategori").Preload("Unit").Preload("HasilAI").Preload("Disposisi").Preload("Disposisi.Unit").Preload("Validasi").Preload("ResponPengaduan.User").
		Order("updated_at DESC").Find(&items).Error
	return items, unitID, err
}

func (s *KasubagService) userUnitID(userID uint64) (uint, error) {
	var user models.User
	if err := config.DB.Select("unit_id").First(&user, userID).Error; err != nil {
		return 0, err
	}
	if user.UnitID == nil || *user.UnitID == 0 {
		return 0, errors.New("akun Kasubag belum terhubung ke unit")
	}
	return *user.UnitID, nil
}

func (s *KasubagService) assertUnit(userID uint64, complaintUnitID *uint) error {
	unitID, err := s.userUnitID(userID)
	if err != nil || complaintUnitID == nil || unitID != *complaintUnitID {
		return ErrForbidden
	}
	return nil
}
