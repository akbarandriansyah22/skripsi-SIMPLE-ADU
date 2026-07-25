package services

import (
	"errors"
	"strings"

	"backend/models"
	"gorm.io/gorm"
)

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func optionalInt64(value int64) *int64 {
	if value <= 0 {
		return nil
	}
	return &value
}

const (
	StatusMenunggu          = "Menunggu Verifikasi"
	StatusDitolak           = "Ditolak"
	StatusMenungguDisposisi = "Menunggu Disposisi"
	StatusDiteruskanUnit    = "Diteruskan ke Unit"
	// StatusDiteruskan is retained for compatibility with old callers.
	StatusDiteruskan = StatusMenungguDisposisi
	StatusDiproses   = "Diproses"
	StatusSelesai    = "Selesai"
)

func isValidStatusTransition(current, next string) bool {
	if current == next {
		return false
	}
	switch strings.ToLower(current) {
	case strings.ToLower(StatusMenunggu):
		return strings.EqualFold(next, StatusDitolak) || strings.EqualFold(next, StatusMenungguDisposisi) || strings.EqualFold(next, StatusDiteruskanUnit)
	case strings.ToLower(StatusMenungguDisposisi):
		return strings.EqualFold(next, StatusDiteruskanUnit)
	case strings.ToLower(StatusDiteruskanUnit):
		return strings.EqualFold(next, StatusDiproses)
	case strings.ToLower(StatusDiproses):
		return strings.EqualFold(next, StatusSelesai)
	default:
		return false
	}
}

func recordStatusChange(tx *gorm.DB, complaintID, changedBy uint, oldStatus, newStatus, note string) error {
	if strings.EqualFold(oldStatus, newStatus) {
		return nil
	}
	if !isValidStatusTransition(oldStatus, newStatus) {
		return errors.New("perubahan status pengaduan tidak valid")
	}
	return tx.Create(&models.RiwayatStatusPengaduan{
		PengaduanID: complaintID,
		ChangedBy:   changedBy,
		StatusLama:  oldStatus,
		StatusBaru:  newStatus,
		Catatan:     strings.TrimSpace(note),
	}).Error
}

func createNotification(tx *gorm.DB, userID, complaintID uint, title, body string) error {
	return tx.Create(&models.Notifikasi{
		UserID:      userID,
		PengaduanID: &complaintID,
		Judul:       title,
		Isi:         body,
	}).Error
}
