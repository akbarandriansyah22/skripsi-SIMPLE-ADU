package migrations

import (
	"embed"

	"gorm.io/gorm"
)

// SQL files are embedded into the binary so Docker does not depend on a host path.
//
//go:embed *.sql
var files embed.FS

func Run(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version varchar(255) PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP)`).Error; err != nil {
		return err
	}
	const version = "001_add_kasubag_and_validation.sql"
	var count int64
	if err := db.Table("schema_migrations").Where("version = ?", version).Count(&count).Error; err != nil {
		return err
	}
	sql, err := files.ReadFile(version)
	if err != nil {
		return err
	}
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Exec(string(sql)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		if err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
