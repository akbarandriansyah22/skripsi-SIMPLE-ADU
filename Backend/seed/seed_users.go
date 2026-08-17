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
	{Nama: "Administrasi", Deskripsi: "Layanan administrasi, surat-menyurat, legalisasi, pembayaran, dan dokumen mahasiswa"},
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
	if err := seedDemoHighUrgencyComplaint(db); err != nil {
		return err
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
			SumberAkun:   "manual",
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
		noHP := strings.TrimSpace(demo.NoHP)
		student := &models.Mahasiswa{UserID: userID, NIM: demo.NIM, ProgramStudi: demo.ProgramStudi, Angkatan: demo.Angkatan}
		if noHP != "" {
			student.NoHP = &noHP
		}
		return tx.Create(student).Error
	}
	return err
}

// seedDemoHighUrgencyComplaint creates one stable complaint for local demos.
// It is intentionally inserted with the same AI result that the local AI
// service would produce for a safety-critical complaint, so the workflow can
// be tested even when the AI container is not available during seeding.
func seedDemoHighUrgencyComplaint(db *gorm.DB) error {
	const (
		demoEmail = "mahasiswa@simpel-adu.test"
		ticket    = "ADU-DEMO-TINGGI-001"
	)

	var user models.User
	if err := db.Where("LOWER(email) = ?", demoEmail).First(&user).Error; err != nil {
		return err
	}

	var category models.KategoriPengaduan
	if err := db.Where("nama = ?", "Fasilitas").First(&category).Error; err != nil {
		return err
	}

	var existing models.Pengaduan
	err := db.Where("kode_tiket = ?", ticket).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	complaint := &models.Pengaduan{
		KodeTiket:  ticket,
		UserID:     user.ID,
		KategoriID: category.ID,
		Judul:      "Korsleting listrik di laboratorium komputer",
		Deskripsi:  "Terjadi korsleting listrik di laboratorium komputer dan muncul asap dari stop kontak. Kondisi ini membahayakan mahasiswa.",
		Status:     "Menunggu Verifikasi",
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(complaint).Error; err != nil {
			return err
		}

		positiveScore := 0
		negativeScore := -4
		result := &models.HasilAI{
			PengaduanID:        complaint.ID,
			CleanedText:        "terjadi korsleting listrik laboratorium komputer muncul asap stop kontak kondisi membahayakan mahasiswa",
			Tokens:             models.JSONB(`["terjadi","korsleting","listrik","laboratorium","komputer","muncul","asap","stop","kontak","kondisi","membahayakan","mahasiswa"]`),
			SkorPositif:        &positiveScore,
			SkorNegatif:        &negativeScore,
			SkorSentimen:       negativeScore,
			Sentimen:           "Negatif",
			PenjelasanSentimen: "Terdapat kata-kata yang menunjukkan keluhan dan kondisi berbahaya.",
			DetailSkor:         models.JSONB(`[]`),
			MatchedWords:       models.JSONB(`[]`),
			UrgencyScore:       3,
			UrgencyReason:      "Terdapat indikator bahaya/keselamatan: asap, korsleting.",
			Urgensi:            "Tinggi",
			DasarUrgensi:       "Terdapat indikator bahaya/keselamatan: asap, korsleting.",
		}
		if err := tx.Create(result).Error; err != nil {
			return err
		}

		var admins []models.User
		if err := tx.Where("role = ? AND is_active = ?", utils.RoleAdminFakultas, true).Find(&admins).Error; err != nil {
			return err
		}
		for _, admin := range admins {
			complaintID := complaint.ID
			if err := tx.Create(&models.Notifikasi{
				UserID:      admin.ID,
				PengaduanID: &complaintID,
				Judul:       "Pengaduan Baru",
				Isi:         "Pengaduan " + complaint.KodeTiket + " menunggu verifikasi.",
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func demoPassword() string {
	if value := strings.TrimSpace(os.Getenv("SEED_DEMO_PASSWORD")); value != "" {
		return value
	}
	return fallbackDemoPassword
}
