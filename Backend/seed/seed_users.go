package seed

import (
	"errors"
	"log"
	"os"
	"strings"

	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

const fallbackDemoPassword = "admin123"

type demoUser struct {
	NamaLengkap  string
	Email        string
	Role         string
	UnitName     string
	NIM          string
	ProgramStudi string
	Angkatan     int
	NoHP         string
}

var demoUsers = []demoUser{
	{NamaLengkap: "Admin Sistem", Email: "admin.sistem@simpel-adu.test", Role: utils.RoleAdminSistem},
	{NamaLengkap: "Admin Fakultas", Email: "admin@simpel-adu.test", Role: utils.RoleAdminFakultas},
	{NamaLengkap: "Pimpinan Fakultas", Email: "pimpinan@simpel-adu.test", Role: utils.RolePimpinan},
	{NamaLengkap: "Kasubag Akademik", Email: "kasubag.akademik@simpel-adu.test", Role: utils.RoleKasubag, UnitName: "Akademik"},
	{NamaLengkap: "Kasubag Sarana dan Prasarana", Email: "kasubag.sarpras@simpel-adu.test", Role: utils.RoleKasubag, UnitName: "Sarana dan Prasarana"},
	{NamaLengkap: "Mahasiswa Demo", Email: "mahasiswa@simpel-adu.test", Role: utils.RoleMahasiswa, NIM: "2026000001", ProgramStudi: "Teknik Informatika", Angkatan: 2026, NoHP: "081234567890"},
}

var demoCategories = []models.KategoriPengaduan{
	{Nama: "Akademik", Deskripsi: "Layanan akademik dan perkuliahan"},
	{Nama: "Fasilitas", Deskripsi: "Sarana dan prasarana fakultas"},
	{Nama: "Kemahasiswaan", Deskripsi: "Layanan kemahasiswaan"},
}

var demoUnits = []models.Unit{
	{NamaUnit: "Akademik", Email: "akademik@simpel-adu.test"},
	{NamaUnit: "Sarana dan Prasarana", Email: "sarpras@simpel-adu.test"},
}

// SeedDemoData is ordered deliberately: categories, units, users, then student profiles.
// It creates missing demo records and refreshes only the known demo passwords;
// active flags and manually assigned units remain unchanged.
func SeedDemoData(db *gorm.DB) error {
	if err := seedCategories(db); err != nil {
		return err
	}
	unitIDs, err := seedUnits(db)
	if err != nil {
		return err
	}
	for _, demo := range demoUsers {
		if err := seedDemoUser(db, demo, unitIDs[demo.UnitName]); err != nil {
			return err
		}
	}
	log.Println("✅ Demo data seeded idempotently")
	return nil
}

// Kept for callers from older code paths.
func SeedDemoUsers(db *gorm.DB) error { return SeedDemoData(db) }

func seedCategories(db *gorm.DB) error {
	for _, category := range demoCategories {
		var existing models.KategoriPengaduan
		err := db.Where("nama = ?", category.Nama).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&category).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func seedUnits(db *gorm.DB) (map[string]uint, error) {
	ids := make(map[string]uint, len(demoUnits))
	for _, unit := range demoUnits {
		var existing models.Unit
		err := db.Where("nama_unit = ?", unit.NamaUnit).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&unit).Error; err != nil {
				return nil, err
			}
			existing = unit
		} else if err != nil {
			return nil, err
		}
		ids[unit.NamaUnit] = existing.ID
	}
	return ids, nil
}

func seedDemoUser(db *gorm.DB, demo demoUser, unitID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		user, err := upsertDemoUser(tx, demo, unitID)
		if err != nil {
			return err
		}
		if demo.Role == utils.RoleMahasiswa {
			return upsertDemoMahasiswa(tx, user.ID, demo)
		}
		return nil
	})
}

func upsertDemoUser(tx *gorm.DB, demo demoUser, unitID uint) (*models.User, error) {
	var user models.User
	email := strings.ToLower(strings.TrimSpace(demo.Email))
	err := tx.Where("LOWER(email) = ?", email).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		passwordHash, hashErr := utils.HashPassword(demoPassword())
		if hashErr != nil {
			return nil, hashErr
		}
		user = models.User{
			NamaLengkap:  demo.NamaLengkap,
			Email:        email,
			PasswordHash: passwordHash,
			Role:         demo.Role,
			IsActive:     true,
		}
		if unitID != 0 {
			user.UnitID = &unitID
		}
		if err := tx.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	// Only the role of a known demo account is normalized. Other user records are untouched.
	passwordHash, hashErr := utils.HashPassword(demoPassword())
	if hashErr != nil {
		return nil, hashErr
	}
	updates := map[string]interface{}{"role": demo.Role, "password_hash": passwordHash}
	if demo.UnitName != "" && user.UnitID == nil && unitID != 0 {
		updates["unit_id"] = unitID
	}
	if err := tx.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	user.Role = demo.Role
	if user.UnitID == nil && unitID != 0 {
		user.UnitID = &unitID
	}
	return &user, nil
}

func upsertDemoMahasiswa(tx *gorm.DB, userID uint, demo demoUser) error {
	var mahasiswa models.Mahasiswa
	err := tx.Where("user_id = ?", userID).First(&mahasiswa).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx.Create(&models.Mahasiswa{UserID: userID, NIM: demo.NIM, ProgramStudi: demo.ProgramStudi, Angkatan: demo.Angkatan, NoHP: demo.NoHP}).Error
	}
	return err
}

func demoPassword() string {
	if value := strings.TrimSpace(os.Getenv("SEED_DEMO_PASSWORD")); value != "" {
		return value
	}
	return fallbackDemoPassword
}
