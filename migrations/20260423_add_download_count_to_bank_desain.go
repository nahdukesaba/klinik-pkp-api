package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Migration20260423AddDownloadCountToBankDesain(db *gorm.DB) error {
	tx := db.Begin()

	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var exists bool

	err := tx.Raw(`SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'bank_design' AND column_name = 'download_count')`).Scan(&exists).Error

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check column existence: %v", err)
	}

	if exists {
		log.Println("!? Column download_count already exists, skipping...")

		return tx.Commit().Error
	}

	// ADD COLUMN WITH DEFAULT VALUE
	if err := tx.Exec(`ALTER TABLE bank_desain ADD COLUMN download_count BIGINT DEFAULT 0`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to add download_count column: %v", err)
	}

	// AUTOMATICALLY BACKFILL NULL VALUES WITH DEFAULT
	if err := tx.Exec(`UPDATE bank_desain SET download_count = 0 WHERE download_count IS NULL`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to backfill download_count column: %v", err)
	}

	// APPLY NOT NULL CONSTRAINT
	if err := tx.Exec(`ALTER TABLE bank_desain ALTER COLUMN download_count SET NOT NULL`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to apply NOT NULL constraint: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Println("✓ Migration completed successfully")

	return nil
}

func Migration20260423AddDownloadCountToBankDesainRollback(db *gorm.DB) error {
	if err := db.Exec(`ALTER TABLE bank_desain DROP COLUMN IF EXISTS download_count`).Error; err != nil {
		return fmt.Errorf("failed to drop status: %v", err)
	}

	log.Println("✓ Migration rolled back successfully")

	return nil
}
