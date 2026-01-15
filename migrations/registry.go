package migrations

import (
	"gorm.io/gorm"
)

// MIGRATION FUNCTION TYPE
type MigrationFunction func(*gorm.DB) error

// MAPS MIGRATION FILENAMES TO THEIR FUNCTIONS
var MigrationRegistry = map[string]MigrationFunction{
	"20260104_add_example_column_to_bsps": Migration20260104AddExampleColumn,
	"20260104_add_indexes_to_bsps":        Migration20260104AddIndexesToBSPS,
}

// MAPS ROLLBACK FILENAMES TO THEIR FUNCTIONS
var RollbackRegistry = map[string]MigrationFunction{
	"20260104_add_example_column_to_bsps": Migration20260104AddExampleColumnRollback,
	"20260104_add_indexes_to_bsps":        Migration20260104AddIndexesToBSPSRollback,
}

// EXECUTE MIGRATION FOR SPECIFIC FILE NAME
func RunModification(db *gorm.DB, name string) error {
	if migration, exists := MigrationRegistry[name]; exists {
		return migration(db)
	}

	return nil
}

// EXECUTE ROLLBACK FOR SPECIFIC FILE NAME
func RunRollback(db *gorm.DB, name string) error {
	if rollback, exists := RollbackRegistry[name]; exists {
		return rollback(db)
	}

	return nil
}

// P.S
// IN if migration, exists := MigrationRegistry[name]; exists {
// 		return migration(db)
// 	}
// THE "exists" IS A SPECIAL GO LANGUAGE FEATURE WHICH A BOOLEAN AUTOMATICALLY RETURNED BY GO