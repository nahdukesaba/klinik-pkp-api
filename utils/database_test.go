package utils

import "testing"

func TestKlinikBrowserSSOMigrationRegistered(t *testing.T) {
	const migration = "20261023_create_klinik_sso_flow"
	if MigrationRegistry[migration] == nil || RollbackRegistry[migration] == nil {
		t.Fatal("Klinik browser SSO migration and rollback must be registered")
	}
}
