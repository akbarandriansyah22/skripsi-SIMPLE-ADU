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

type AdminSistemService struct{}

func NewAdminSistemService() *AdminSistemService { return &AdminSistemService{} }

func (s *AdminSistemService) Dashboard() (map[string]int64, error) {
	var users, complaints, units, categories int64
	if err := config.DB.Model(&models.User{}).Count(&users).Error; err != nil {
		return nil, err
	}
	if err := config.DB.Model(&models.Pengaduan{}).Count(&complaints).Error; err != nil {
		return nil, err
	}
	if err := config.DB.Model(&models.Unit{}).Count(&units).Error; err != nil {
		return nil, err
	}
	if err := config.DB.Model(&models.KategoriPengaduan{}).Count(&categories).Error; err != nil {
		return nil, err
	}
	return map[string]int64{"total_users": users, "total_pengaduan": complaints, "total_units": units, "total_categories": categories}, nil
}

func (s *AdminSistemService) Users() ([]models.User, error) {
	var users []models.User
	err := config.DB.Preload("Unit").Order("created_at DESC").Find(&users).Error
	return users, err
}

func (s *AdminSistemService) User(id uint64) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Unit").First(&user, id).Error
	return &user, err
}

func canonicalAccount(role string, unitID *uint) (string, *uint, error) {
	canonical := utils.CanonicalRole(role)
	if canonical == "" {
		return "", nil, errors.New("role tidak valid")
	}
	if canonical == utils.RoleKasubag {
		if unitID == nil || *unitID == 0 {
			return "", nil, errors.New("role kasubag wajib memiliki unit_id")
		}
		var unit models.Unit
		if err := config.DB.First(&unit, *unitID).Error; err != nil {
			return "", nil, errors.New("unit tidak ditemukan")
		}
		return canonical, unitID, nil
	}
	return canonical, nil, nil
}

func (s *AdminSistemService) CreateUser(req dto.CreateUserRequest) (*models.User, error) {
	role, unitID, err := canonicalAccount(req.Role, req.UnitID)
	if err != nil {
		return nil, err
	}
	returnUser := &models.User{}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.User{}).Where("LOWER(email) = ?", strings.ToLower(strings.TrimSpace(req.Email))).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("email sudah digunakan")
		}
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		returnUser = &models.User{NamaLengkap: strings.TrimSpace(req.NamaLengkap), Email: strings.ToLower(strings.TrimSpace(req.Email)), PasswordHash: hash, Role: role, UnitID: unitID, IsActive: true, PasswordMustChange: false}
		return tx.Create(returnUser).Error
	})
	return returnUser, err
}

func (s *AdminSistemService) UpdateUser(id uint64, req dto.UpdateUserRequest) (*models.User, error) {
	role, unitID, err := canonicalAccount(req.Role, req.UnitID)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&models.User{}).Where("LOWER(email) = ? AND id <> ?", strings.ToLower(strings.TrimSpace(req.Email)), id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("email sudah digunakan")
		}
		return tx.Model(&user).Updates(map[string]interface{}{"nama_lengkap": strings.TrimSpace(req.NamaLengkap), "email": strings.ToLower(strings.TrimSpace(req.Email)), "role": role, "unit_id": unitID}).Error
	})
	if err != nil {
		return nil, err
	}
	return s.User(id)
}

func (s *AdminSistemService) SetUserStatus(id uint64, active bool) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).Update("is_active", active).Error
}

func (s *AdminSistemService) ResetPassword(id uint64, password string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return config.DB.Model(&models.User{}).Where("id = ?", id).Updates(map[string]any{"password_hash": hash, "password_must_change": false}).Error
}

func (s *AdminSistemService) Units() ([]models.Unit, error) {
	var rows []models.Unit
	err := config.DB.Order("nama_unit ASC").Find(&rows).Error
	return rows, err
}
func (s *AdminSistemService) CreateUnit(req dto.UnitRequest) (*models.Unit, error) {
	row := &models.Unit{NamaUnit: strings.TrimSpace(req.NamaUnit), Email: strings.TrimSpace(req.Email), IsActive: true}
	return row, config.DB.Create(row).Error
}
func (s *AdminSistemService) UpdateUnit(id uint64, req dto.UnitRequest) (*models.Unit, error) {
	var row models.Unit
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	err := config.DB.Model(&row).Updates(map[string]interface{}{"nama_unit": strings.TrimSpace(req.NamaUnit), "email": strings.TrimSpace(req.Email)}).Error
	return &row, err
}
func (s *AdminSistemService) SetUnitStatus(id uint64, active bool) error {
	result := config.DB.Model(&models.Unit{}).Where("id = ?", id).Update("is_active", active)
	return result.Error
}
func (s *AdminSistemService) Categories() ([]models.KategoriPengaduan, error) {
	var rows []models.KategoriPengaduan
	err := config.DB.Order("nama ASC").Find(&rows).Error
	return rows, err
}
func (s *AdminSistemService) CreateCategory(req dto.CategoryRequest) (*models.KategoriPengaduan, error) {
	row := &models.KategoriPengaduan{Nama: strings.TrimSpace(req.Nama), Deskripsi: strings.TrimSpace(req.Deskripsi), IsActive: true}
	return row, config.DB.Create(row).Error
}
func (s *AdminSistemService) UpdateCategory(id uint64, req dto.CategoryRequest) (*models.KategoriPengaduan, error) {
	var row models.KategoriPengaduan
	if err := config.DB.First(&row, id).Error; err != nil {
		return nil, err
	}
	err := config.DB.Model(&row).Updates(map[string]interface{}{"nama": strings.TrimSpace(req.Nama), "deskripsi": strings.TrimSpace(req.Deskripsi)}).Error
	return &row, err
}
func (s *AdminSistemService) SetCategoryStatus(id uint64, active bool) error {
	result := config.DB.Model(&models.KategoriPengaduan{}).Where("id = ?", id).Update("is_active", active)
	return result.Error
}
