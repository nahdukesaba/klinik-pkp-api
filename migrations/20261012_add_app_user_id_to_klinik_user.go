package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Migration20261012AddAppUserIDToKlinikUser(db *gorm.DB) error {
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

	// 1. ADD app_user_id COLUMN IF IT DOES NOT EXIST
	var columnExists bool

	if err := tx.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'public'
			  AND table_name = 'klinikpkp_user'
			  AND column_name = 'app_user_id'
		)
	`).Scan(&columnExists).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check app_user_id column: %v", err)
	}

	if !columnExists {
		if err := tx.Exec(`
			ALTER TABLE public.klinikpkp_user
			ADD COLUMN app_user_id UUID
		`).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add app_user_id column: %v", err)
		}

		log.Println("✓ app_user_id column added")
	} else {
		log.Println("!? app_user_id column already exists, continuing...")
	}

	// 2. BACKFILL app_user_id FROM CENTRAL app_users USING NORMALIZED EMAIL
	if err := tx.Exec(`
		UPDATE public.klinikpkp_user AS ku
		SET app_user_id = au.id
		FROM public.app_users AS au
		WHERE ku.app_user_id IS NULL
		  AND LOWER(BTRIM(ku.email)) = LOWER(BTRIM(au.email))
		  AND au.deleted_at IS NULL
	`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to backfill app_user_id: %v", err)
	}

	// 3. VERIFY THAT EVERY ACTIVE, NON-DELETED KLINIK USER IS LINKED
	var unlinkedActiveUsers int64

	if err := tx.Raw(`
		SELECT COUNT(*)
		FROM public.klinikpkp_user
		WHERE deleted_at IS NULL
		  AND is_active = TRUE
		  AND app_user_id IS NULL
	`).Scan(&unlinkedActiveUsers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to verify active Klinik users: %v", err)
	}

	if unlinkedActiveUsers > 0 {
		tx.Rollback()
		return fmt.Errorf(
			"migration aborted: %d active Klinik user(s) are not linked to app_users",
			unlinkedActiveUsers,
		)
	}

	// 4. ADD FOREIGN KEY TO app_users(id) IF IT DOES NOT EXIST
	var foreignKeyExists bool

	if err := tx.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_constraint
			WHERE conname = 'fk_klinikpkp_user_app_user'
			  AND conrelid = 'public.klinikpkp_user'::regclass
		)
	`).Scan(&foreignKeyExists).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check app_user_id foreign key: %v", err)
	}

	if !foreignKeyExists {
		if err := tx.Exec(`
			ALTER TABLE public.klinikpkp_user
			ADD CONSTRAINT fk_klinikpkp_user_app_user
			FOREIGN KEY (app_user_id)
			REFERENCES public.app_users(id)
		`).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add app_user_id foreign key: %v", err)
		}

		log.Println("✓ app_user_id foreign key added")
	}

	// 5. ADD UNIQUE CONSTRAINT SO ONE CENTRAL USER MAPS TO AT MOST ONE KLINIK USER
	var uniqueConstraintExists bool

	if err := tx.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_constraint
			WHERE conname = 'uq_klinikpkp_user_app_user'
			  AND conrelid = 'public.klinikpkp_user'::regclass
		)
	`).Scan(&uniqueConstraintExists).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check app_user_id unique constraint: %v", err)
	}

	if !uniqueConstraintExists {
		if err := tx.Exec(`
			ALTER TABLE public.klinikpkp_user
			ADD CONSTRAINT uq_klinikpkp_user_app_user
			UNIQUE (app_user_id)
		`).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add app_user_id unique constraint: %v", err)
		}

		log.Println("✓ app_user_id unique constraint added")
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Println("✓ Klinik app_user_id migration completed successfully")

	return nil
}

func Migration20261012AddAppUserIDToKlinikUserRollback(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start rollback transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Exec(`
		ALTER TABLE public.klinikpkp_user
		DROP CONSTRAINT IF EXISTS uq_klinikpkp_user_app_user
	`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to drop app_user_id unique constraint: %v", err)
	}

	if err := tx.Exec(`
		ALTER TABLE public.klinikpkp_user
		DROP CONSTRAINT IF EXISTS fk_klinikpkp_user_app_user
	`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to drop app_user_id foreign key: %v", err)
	}

	if err := tx.Exec(`
		ALTER TABLE public.klinikpkp_user
		DROP COLUMN IF EXISTS app_user_id
	`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to drop app_user_id column: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %v", err)
	}

	log.Println("✓ Klinik app_user_id migration rolled back successfully")

	return nil
}
