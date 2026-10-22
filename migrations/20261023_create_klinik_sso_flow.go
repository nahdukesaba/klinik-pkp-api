package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Migration20261023CreateKlinikSSOFlow(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %v", tx.Error)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			tx.Rollback()
			panic(recovered)
		}
	}()

	if err := tx.Exec(`
CREATE TABLE IF NOT EXISTS public.klinikpkp_sso_flow (
    state_hash BYTEA PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ NULL
)
`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create klinikpkp_sso_flow: %v", err)
	}
	if err := tx.Exec(`
CREATE INDEX IF NOT EXISTS idx_klinikpkp_sso_flow_active_expiry
ON public.klinikpkp_sso_flow (expires_at)
WHERE used_at IS NULL
`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to index klinikpkp_sso_flow: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}
	log.Println("Klinik browser SSO flow migration completed successfully")
	return nil
}

func Migration20261023CreateKlinikSSOFlowRollback(db *gorm.DB) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start rollback transaction: %v", tx.Error)
	}
	if err := tx.Exec(`DROP TABLE IF EXISTS public.klinikpkp_sso_flow`).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to drop klinikpkp_sso_flow: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %v", err)
	}
	log.Println("Klinik browser SSO flow migration rolled back successfully")
	return nil
}
