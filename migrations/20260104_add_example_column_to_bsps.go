package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

// EXAMPLE MIGRATION TO ADD A NEW COLUMN TO penerimaan_bsps TABLE
func Migration20260104AddExampleColumn(db *gorm.DB) error {
	log.Println("Running migration: 20260104_add_example_column_to_bsps")

	var exists bool

	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'penerimaan_bsps' 
			AND column_name = 'example_column'
		)
	`).Scan(&exists).Error

	if err != nil {
		return fmt.Errorf("failed to check if column exists: %v", err)
	}

	if exists {
		log.Println("Column example_column already exists, skipping...")
		return nil
	}

	if err := db.Exec(`
		ALTER TABLE penerimaan_bsps 
		ADD COLUMN example_column VARCHAR(255) DEFAULT NULL
	`).Error; err != nil {
		return fmt.Errorf("failed to add example_column: %v", err)
	}

	log.Println("✓ Successfully added example_column to penerimaan_bsps")
	return nil
}

// EXAMPLE ROLLBACK TO REMOVE THE NEW COLUMN FROM penerimaan_bsps TABLE
func Migration20260104AddExampleColumnRollback(db *gorm.DB) error {
	log.Println("Rolling back migration: 20260104_add_example_column_to_bsps")

	if err := db.Exec(`
		ALTER TABLE penerimaan_bsps 
		DROP COLUMN IF EXISTS example_column
	`).Error; err != nil {
		return fmt.Errorf("failed to drop example_column: %v", err)
	}

	log.Println("✓ Successfully removed example_column from penerimaan_bsps")
	return nil
}