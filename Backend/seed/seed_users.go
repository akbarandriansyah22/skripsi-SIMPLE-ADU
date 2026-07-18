package seed

import (
	"errors"
	"log"
	"strings"

	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

const demoPassword = "Admin123!"

type demoUser struct {
	NamaLengkap  string
	Email        string
	Role         string
	NIM          string
	ProgramStudi string
	Angkatan     int
	NoHP         string
}

var demoUsers = []demoUser{
	{
		NamaLengkap: "Admin Fakultas",
		Email:       "admin@simpel-adu.test",
		Role:        utils.RoleAdmin,
	},
	{
		NamaLengkap: "Pimpinan Fakultas",
		Email:       "pimpinan@simpel-adu.test",
		Role:        utils.RolePimpinan,
	},
	{
		NamaLengkap:  "Mahasiswa Demo",
		Email:        "mahasiswa@simpel-adu.test",
		Role:         "mahasiswa",
		NIM:          "2026000001",
		ProgramStudi: "Teknik Informatika",
		Angkatan:     2026,
		NoHP:         "081234567890",
	},
}

func SeedDemoUsers(db *gorm.DB) error {
	for _, demo := range demoUsers {
		if err := seedDemoUser(db, demo); err != nil {
			return err
		}
	}

	log.Println("✅ Demo users seeded")
	return nil
}

func seedDemoUser(db *gorm.DB, demo demoUser) error {
	return db.Transaction(func(tx *gorm.DB) error {
		user, err := upsertDemoUser(tx, demo)
		if err != nil {
			return err
		}

		if demo.Role != "mahasiswa" {
			return nil
		}

		return upsertDemoMahasiswa(tx, user.ID, demo)
	})
}

func upsertDemoUser(tx *gorm.DB, demo demoUser) (*models.User, error) {
	var user models.User
	email := strings.ToLower(strings.TrimSpace(demo.Email))
	err := tx.Where("LOWER(email) = ?", email).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		passwordHash, hashErr := utils.HashPassword(demoPassword)
		if hashErr != nil {
			return nil, hashErr
		}
		user.NamaLengkap = demo.NamaLengkap
		user.Email = email
		user.PasswordHash = passwordHash
		user.Role = demo.Role
		user.IsActive = true
		if err := tx.Create(&user).Error; err != nil {
			return nil, err
		}

		return &user, nil
	}

	// Existing accounts are never reset or reactivated by a startup seeder.
	// Legacy demo roles are normalized so JWT and route authorization converge.
	updates := map[string]interface{}{"role": demo.Role}
	if err := tx.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func upsertDemoMahasiswa(tx *gorm.DB, userID uint, demo demoUser) error {
	var mahasiswa models.Mahasiswa
	err := tx.Where("user_id = ?", userID).First(&mahasiswa).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	mahasiswa.UserID = userID
	mahasiswa.NIM = demo.NIM
	mahasiswa.ProgramStudi = demo.ProgramStudi
	mahasiswa.Angkatan = demo.Angkatan
	mahasiswa.NoHP = demo.NoHP

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx.Create(&mahasiswa).Error
	}

	return tx.Save(&mahasiswa).Error
}
