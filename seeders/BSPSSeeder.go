package seeders

import (
	"klinik-pkp-api/internal/api/bsps"
	"log"

	"gorm.io/gorm"
)

// SEEDS FROM PREDEFINED DATA
func BSPSSeeder(db *gorm.DB) error {
	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&bsps.BSPS{}).Count(&count)

	if count > 0 {
		log.Println("!? BSPS table already seeded, skipping...")
		return nil
	}

	// TODO: Implement BSPS seeding logic if needed
	log.Println("✓ BSPS table seeder executed (no data to seed)")

	return nil
}
