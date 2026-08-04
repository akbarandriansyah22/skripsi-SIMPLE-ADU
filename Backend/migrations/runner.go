package migrations

import (
	"embed"
	"fmt"
	"sort"

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
	entries, err := files.ReadDir(".")
	if err != nil {
		return err
	}
	versions := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name()[len(entry.Name())-4:] == ".sql" {
			versions = append(versions, entry.Name())
		}
	}
	sort.Strings(versions)
	for _, version := range versions {
		var count int64
		if err := db.Table("schema_migrations").Where("version = ?", version).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
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
			return fmt.Errorf("migration %s: %w", version, err)
		}
		if err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			return err
		}
	}
	return nil
}
