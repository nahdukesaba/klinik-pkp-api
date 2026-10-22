package utils

import "klinik-pkp-api/migrations"

func init() {
	MigrationRegistry["20261023_create_klinik_sso_flow"] =
		migrations.Migration20261023CreateKlinikSSOFlow
	RollbackRegistry["20261023_create_klinik_sso_flow"] =
		migrations.Migration20261023CreateKlinikSSOFlowRollback
}
