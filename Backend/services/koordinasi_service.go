package services

import (
	"errors"
	"strings"

	dto "backend/DTO"
	"backend/config"
	"backend/models"
	"backend/utils"
	"gorm.io/gorm"
)

type KoordinasiService struct{}

func NewKoordinasiService() *KoordinasiService { return &KoordinasiService{} }

func (s *KoordinasiService) assertAccess(userID uint, role string, complaintID uint) (*models.Pengaduan, error) {
	var complaint models.Pengaduan
	if err := config.DB.Preload("HasilAI").First(&complaint, complaintID).Error; err != nil {
		return nil, err
	}
	if complaint.HasilAI == nil || !strings.EqualFold(complaint.HasilAI.Urgensi, "Tinggi") {
		return nil, errors.New("koordinasi internal hanya tersedia untuk urgensi Tinggi")
	}
	switch utils.CanonicalRole(role) {
	case utils.RolePimpinan, utils.RoleAdminFakultas:
		return &complaint, nil
	case utils.RoleKasubag:
		if !userBelongsToComplaintUnit(uint64(userID), complaint.UnitID) {
			return nil, ErrForbidden
		}
		return &complaint, nil
	default:
		return nil, ErrForbidden
	}
}

func (s *KoordinasiService) List(userID uint, role string, complaintID uint) ([]dto.KoordinasiResponse, error) {
	if _, err := s.assertAccess(userID, role, complaintID); err != nil {
		return nil, err
	}
	var rows []models.KoordinasiInternal
	if err := config.DB.Preload("Sender").Where("pengaduan_id = ?", complaintID).Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]dto.KoordinasiResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapCoordination(row))
	}
	return result, nil
}

func (s *KoordinasiService) Create(userID uint, role string, complaintID uint, parentID *uint, message, attachment, original, mime string, size int64) error {
	complaint, err := s.assertAccess(userID, role, complaintID)
	if err != nil {
		return err
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return errors.New("pesan koordinasi wajib diisi")
	}
	if parentID != nil {
		var parent models.KoordinasiInternal
		if err := config.DB.Where("id = ? AND pengaduan_id = ?", *parentID, complaintID).First(&parent).Error; err != nil {
			return errors.New("pesan yang dibalas tidak ditemukan")
		}
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		row := &models.KoordinasiInternal{PengaduanID: complaint.ID, SenderID: userID, ParentID: parentID, Pesan: message, Lampiran: attachment, LampiranNamaAsli: optionalString(original), LampiranMimeType: optionalString(mime), LampiranUkuran: optionalInt64(size)}
		if err := tx.Create(row).Error; err != nil {
			return err
		}
		var recipients []models.User
		canonical := utils.CanonicalRole(role)
		if canonical == utils.RolePimpinan || canonical == utils.RoleAdminFakultas {
			if complaint.UnitID != nil {
				if err := tx.Where("role = ? AND unit_id = ? AND is_active = ?", utils.RoleKasubag, *complaint.UnitID, true).Find(&recipients).Error; err != nil {
					return err
				}
			}
		} else {
			if err := tx.Where("role = ? AND is_active = ?", utils.RolePimpinan, true).Find(&recipients).Error; err != nil {
				return err
			}
		}
		for _, recipient := range recipients {
			if recipient.ID != userID {
				if err := createNotification(tx, recipient.ID, complaint.ID, "Koordinasi Internal Baru", "Ada pesan koordinasi baru untuk pengaduan "+complaint.KodeTiket+"."); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func mapCoordination(row models.KoordinasiInternal) dto.KoordinasiResponse {
	role := utils.CanonicalRole(row.Sender.Role)
	return dto.KoordinasiResponse{ID: row.ID, PengaduanID: row.PengaduanID, ParentID: row.ParentID, SenderID: row.SenderID, SenderName: row.Sender.NamaLengkap, SenderRole: role, Pesan: row.Pesan, Lampiran: row.Lampiran, LampiranNamaAsli: stringValue(row.LampiranNamaAsli), CreatedAt: row.CreatedAt}
}
