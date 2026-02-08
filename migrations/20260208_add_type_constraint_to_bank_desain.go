package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Migration20260208AddTypeConstraintToBankDesain(db *gorm.DB) error {
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

	// CHECK IF CONSTRAINT ALREADY EXISTS
	err := tx.Raw(`SELECT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE table_schema = 'public' AND table_name = 'bank_desain' AND constraint_name = 'chk_bank_desain_type')`).Scan(&exists).Error

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check constraint existence: %v", err)
	}

	if exists {
		log.Println("!? Constraint chk_bank_desain_type already exists, skipping...")

		return tx.Commit().Error
	}

	// CHECK IF VALUE IS WITHIN ALLOWED SET
	if err := tx.Exec(`ALTER TABLE bank_desain ADD CONSTRAINT chk_bank_desain_type CHECK (type IN ('Tipe 36', 'Tipe 45', 'Tipe 54', 'Rusun'))`).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to add CHECK constraint: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Println("✓ Migration completed successfully")

	return nil
}

// ROLLBACK TO REMOVE THE CONSTRAINT FROM bank_desain TABLE
func Migration20260208AddTypeConstraintToBankDesainRollback(db *gorm.DB) error {
	if err := db.Exec(`ALTER TABLE bank_desain DROP CONSTRAINT IF EXISTS chk_bank_desain_type`).Error; err != nil {
		return fmt.Errorf("failed to drop constraint: %v", err)
	}

	log.Println("✓ Migration rolled back successfully")

	return nil
}
