package services

import (
	dto "backend/DTO"
	"backend/models"
	"backend/repository"
)

type NotifikasiService struct {
	repo repository.NotifikasiRepository
}

func NewNotifikasiService() *NotifikasiService {
	return &NotifikasiService{
		repo: repository.NewNotifikasiRepository(),
	}
}

func (s *NotifikasiService) Create(userID uint, judul string, isi string) error {
	return s.repo.Create(&models.Notifikasi{
		UserID: userID,
		Judul:  judul,
		Isi:    isi,
	})
}

func (s *NotifikasiService) GetByUserID(userID uint64) ([]dto.NotifikasiResponse, error) {
	items, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return mapNotifikasiResponses(items), nil
}

func (s *NotifikasiService) GetUnreadByUserID(userID uint64) ([]dto.NotifikasiResponse, error) {
	items, err := s.repo.GetUnreadByUserID(userID)
	if err != nil {
		return nil, err
	}

	return mapNotifikasiResponses(items), nil
}

func (s *NotifikasiService) MarkAsRead(userID uint64, id uint64) error {
	return s.repo.MarkAsReadByUserID(id, userID)
}

func mapNotifikasiResponses(items []models.Notifikasi) []dto.NotifikasiResponse {
	responses := make([]dto.NotifikasiResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, dto.NotifikasiResponse{
			ID:          item.ID,
			PengaduanID: item.PengaduanID,
			Judul:       item.Judul,
			Isi:         item.Isi,
			IsRead:      item.IsRead,
			CreatedAt:   item.CreatedAt,
			ReadAt:      item.ReadAt,
		})
	}

	return responses
}
