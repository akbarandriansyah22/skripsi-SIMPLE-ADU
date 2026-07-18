package config

import (
	"fmt"
	"log"
	"os"

	"backend/models"
	"backend/seed"
	"backend/utils"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	_ = godotenv.Load()
	if err := utils.ValidateJWTSecret(); err != nil {
		log.Fatalf("Konfigurasi JWT tidak valid: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	DB = db

	log.Println("✅ Database Connected")

	if err := autoMigrate(); err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}

	if err := seed.SeedDemoUsers(DB); err != nil {
		log.Fatalf("Seeder demo user gagal: %v", err)
	}
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Mahasiswa{},
		&models.KategoriPengaduan{},
		&models.Unit{},
		&models.Pengaduan{},
		&models.HasilAI{},
		&models.Disposisi{},
		&models.ResponPengaduan{},
		&models.Notifikasi{},
	)
}
