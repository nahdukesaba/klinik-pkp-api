package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"log"
)

// Migration20260104AddIndexesToBSPS adds indexes to improve query performance
func Migration20260104AddIndexesToBSPS(db *gorm.DB) error {
	log.Println("📝 Running migration: 20260104_add_indexes_to_bsps")

	// Add index on village_id
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_penerimaan_bsps_village_id 
		ON penerimaan_bsps(village_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create index on village_id: %v", err)
	}
	log.Println("✓ Created index on village_id")

	// Add index on district_id
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_penerimaan_bsps_district_id 
		ON penerimaan_bsps(district_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create index on district_id: %v", err)
	}
	log.Println("✓ Created index on district_id")

	// Add index on region_id
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_penerimaan_bsps_region_id 
		ON penerimaan_bsps(region_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create index on region_id: %v", err)
	}
	log.Println("✓ Created index on region_id")

	// Add composite index on year_given and region_id for common queries
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_penerimaan_bsps_year_region 
		ON penerimaan_bsps(year_given, region_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create composite index: %v", err)
	}
	log.Println("✓ Created composite index on year_given and region_id")

	log.Println("✓ Successfully added all indexes to penerimaan_bsps")
	return nil
}

// Migration20260104AddIndexesToBSPSRollback removes the indexes
func Migration20260104AddIndexesToBSPSRollback(db *gorm.DB) error {
	log.Println("⏮️  Rolling back migration: 20260104_add_indexes_to_bsps")

	indexes := []string{
		"idx_penerimaan_bsps_village_id",
		"idx_penerimaan_bsps_district_id",
		"idx_penerimaan_bsps_region_id",
		"idx_penerimaan_bsps_year_region",
	}

	for _, idx := range indexes {
		if err := db.Exec(fmt.Sprintf("DROP INDEX IF EXISTS %s", idx)).Error; err != nil {
			log.Printf("⚠️  Warning: failed to drop index %s: %v", idx, err)
		} else {
			log.Printf("✓ Dropped index %s", idx)
		}
	}

	log.Println("✓ Successfully removed all indexes from penerimaan_bsps")
	return nil
}
