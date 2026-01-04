package utils

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func DropTable(db *gorm.DB, table string) error {
	log.Printf("✓ %s table dropped", cases.Title(language.Und).String(table))

	if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
		return fmt.Errorf("failed to drop table %s: %v", table, err)
	}

	return nil
}

func TruncateTable(db *gorm.DB, table string) error {
	log.Printf("✓ %s table truncated", cases.Title(language.Und).String(table))

	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
		return fmt.Errorf("failed to truncate table %s: %v", table, err)
	}

	return nil
}
