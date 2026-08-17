package services

import (
	"errors"
	"strings"
	"time"

	dto "backend/DTO"
	"backend/models"
	"backend/repository"
	"backend/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo: repository.NewAuthRepository(),
	}
}

// =======================================
// REGISTER
// =======================================

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.LoginResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.NIM = strings.TrimSpace(req.NIM)

	// Cek Email
	_, err := s.repo.GetUserByEmail(req.Email)

	if err == nil {
		return nil, errors.New("email sudah digunakan")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Cek NIM
	_, err = s.repo.GetMahasiswaByNIM(req.NIM)

	if err == nil {
		return nil, errors.New("NIM sudah terdaftar")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash Password
	hashPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	// User
	user := &models.User{
		NamaLengkap:  req.NamaLengkap,
		Email:        req.Email,
		PasswordHash: hashPassword,
		Role:         "mahasiswa",
		IsActive:     true,
		SumberAkun:   "manual",
	}

	// Mahasiswa
	mahasiswa := &models.Mahasiswa{
		NIM:          req.NIM,
		ProgramStudi: req.ProgramStudi,
		Angkatan:     req.Angkatan,
		NoHP:         optionalString(req.NoHP),
	}

	// Simpan ke database
	if err := s.repo.RegisterMahasiswa(user, mahasiswa); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID,
			NamaLengkap: user.NamaLengkap,
			Email:       user.Email,
			Role:        user.Role,
			IsActive:    user.IsActive,
			UnitID:      user.UnitID,
			UnitName:    userUnitName(*user),
		},
	}, nil
}

// =======================================
// LOGIN
// =======================================

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {

	identifier := strings.TrimSpace(req.Email)
	if strings.Contains(identifier, "@") {
		identifier = strings.ToLower(identifier)
	}
	user, err := s.repo.GetUserByEmailOrNIM(identifier)

	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Akun tidak aktif
	if !user.IsActive {
		return nil, errors.New("akun tidak aktif")
	}

	// Cek Password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("email atau password salah")
	}

	role := utils.CanonicalRole(user.Role)
	if role == "" {
		return nil, errors.New("role akun tidak valid")
	}
	token, err := utils.GenerateJWT(user.ID, role)

	if err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.repo.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("last_login_at", &now).Error

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:                 user.ID,
			NamaLengkap:        user.NamaLengkap,
			Email:              user.Email,
			Role:               role,
			IsActive:           user.IsActive,
			PasswordMustChange: user.PasswordMustChange,
			UnitID:             user.UnitID,
			UnitName:           userUnitName(*user),
		},
	}, nil
}

func (s *AuthService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	var user models.User
	if err := s.repo.DB.First(&user, userID).Error; err != nil {
		return err
	}
	if req.CurrentPassword != "" && !utils.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return errors.New("password saat ini salah")
	}
	hash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	return s.repo.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{"password_hash": hash, "must_change_password": false}).Error
}

// =======================================
// PROFILE
// =======================================

func (s *AuthService) Profile(userID uint) (*dto.UserResponse, error) {

	user, err := s.repo.GetUserByID(userID)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:                 user.ID,
		NamaLengkap:        user.NamaLengkap,
		Email:              user.Email,
		Role:               utils.CanonicalRole(user.Role),
		IsActive:           user.IsActive,
		PasswordMustChange: user.PasswordMustChange,
		UnitID:             user.UnitID,
		UnitName:           userUnitName(*user),
	}, nil
}

func userUnitName(user models.User) string {
	if user.Unit == nil {
		return ""
	}
	return user.Unit.NamaUnit
}
