package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

func Migration20260201AddStatusColumnToBSPS(db *gorm.DB) error {
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

	err := tx.Raw(`SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'penerimaan_bsps' AND column_name = 'status')`).Scan(&exists).Error

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check column existence: %v", err)
	}

	if exists {
		log.Println("!? Column status already exists, skipping...")

		return tx.Commit().Error
	}

	// ADD COLUMN WITH DEFAULT VALUE
	if err := tx.Exec(`ALTER TABLE penerimaan_bsps ADD COLUMN status VARCHAR(255) DEFAULT 'Rencana'`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to add status column: %v", err)
	}

	// AUTOMATICALLY BACKFILL NULL VALUES WITH DEFAULT
	if err := tx.Exec(`UPDATE penerimaan_bsps SET status = 'Rencana' WHERE status IS NULL`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to backfill status column: %v", err)
	}

	// APPLY NOT NULL CONSTRAINT
	if err := tx.Exec(`ALTER TABLE penerimaan_bsps ALTER COLUMN status SET NOT NULL`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to apply NOT NULL constraint: %v", err)
	}

	// CHECK IF VALUE IS WITHIN ALLOWED SET
	if err := tx.Exec(`ALTER TABLE penerimaan_bsps ADD CONSTRAINT chk_penerimaan_bsps_status CHECK (status IN ('Rencana', 'Dalam Proses', 'Selesai'))`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to add CHECK constraint: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Println("✓ Migration completed successfully")

	return nil
}

// EXAMPLE ROLLBACK TO REMOVE THE NEW COLUMN FROM penerimaan_bsps TABLE
func Migration20260201AddStatusColumnToBSPSRollback(db *gorm.DB) error {
	if err := db.Exec(`ALTER TABLE penerimaan_bsps DROP COLUMN IF EXISTS status`).Error; err != nil {
		return fmt.Errorf("failed to drop status: %v", err)
	}

	log.Println("✓ Migration rolled back successfully")

	return nil
}
