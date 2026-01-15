package utils

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"log"
)

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
	if err := db.Exec(fmt.Sprintf("DROP TABLE %s CASCADE", table)).Error; err != nil {
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
	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
		return fmt.Errorf("failed to truncate table %s: %v", table, err)
	}

	log.Printf("✓ %s table truncated", cases.Title(language.Und).String(table))

	return nil
}
