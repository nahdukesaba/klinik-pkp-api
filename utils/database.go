package utils

import (
	"fmt"
	"log"

	"klinik-pkp-api/migrations"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// MIGRATION FUNCTION TYPE
type MigrationFunction func(*gorm.DB) error

// MAPS MIGRATION FILENAMES TO THEIR FUNCTIONS
var MigrationRegistry = map[string]MigrationFunction{
	"20260201_add_status_column_to_bsps":          migrations.Migration20260201AddStatusColumnToBSPS,
	"20260208_add_type_constraint_to_bank_desain": migrations.Migration20260208AddTypeConstraintToBankDesain,
	"20260423_add_download_count_to_bank_desain":  migrations.Migration20260423AddDownloadCountToBankDesain,
}

// MAPS ROLLBACK FILENAMES TO THEIR FUNCTIONS
var RollbackRegistry = map[string]MigrationFunction{
	"20260201_add_status_column_to_bsps":          migrations.Migration20260201AddStatusColumnToBSPSRollback,
	"20260208_add_type_constraint_to_bank_desain": migrations.Migration20260208AddTypeConstraintToBankDesainRollback,
	"20260423_add_download_count_to_bank_desain":  migrations.Migration20260423AddDownloadCountToBankDesainRollback,
}

func DropTable(db *gorm.DB, table string) error {
	// CHECK IF TABLE EXISTS FIRST
	var exists bool

	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = ?
		)
	`, table).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check if table %s exists: %v", table, err)
	}

	if !exists {
		return fmt.Errorf("table %s does not exist", table)
	}

	// DROP THE TABLE
	if err := db.Exec(fmt.Sprintf(`DROP TABLE "%s" CASCADE`, table)).Error; err != nil {
		return fmt.Errorf("failed to drop table %s: %v", table, err)
	}

	log.Printf("✓ %s table dropped", cases.Title(language.Und).String(table))

	return nil
}

func TruncateTable(db *gorm.DB, table string) error {
	// CHECK IF TABLE EXISTS FIRST
	var exists bool

	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = ?
		)
	`, table).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check if table %s exists: %v", table, err)
	}

	if !exists {
		return fmt.Errorf("table %s does not exist", table)
	}

	// TRUNCATE THE TABLE
	if err := db.Exec(fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY CASCADE`, table)).Error; err != nil {
		return fmt.Errorf("failed to truncate table %s: %v", table, err)
	}

	log.Printf("✓ %s table truncated", cases.Title(language.Und).String(table))

	return nil
}

// EXECUTE MIGRATION FOR SPECIFIC FILE NAME
func RunModification(db *gorm.DB, filename string) error {
	if migration, exists := MigrationRegistry[filename]; exists {
		return migration(db)
	}

	return fmt.Errorf("migration %s not found", filename)
}

// EXECUTE ROLLBACK FOR SPECIFIC FILE NAME
func RunRollback(db *gorm.DB, filename string) error {
	if rollback, exists := RollbackRegistry[filename]; exists {
		return rollback(db)
	}

	return nil
}

// P.S
// IN if migration, exists := MigrationRegistry[name]; exists {
// 		return migration(db)
// 	}
// THE "exists" IS A SPECIAL GO LANGUAGE FEATURE WHICH DATA TYPE IS BOOLEAN AND RETURNED AUTOMATICALLY BY GO
