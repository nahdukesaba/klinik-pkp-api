package seed

import (
	"klinik-api/config"
	"log"

	"gorm.io/gorm"
)

// Run menjalankan seeder secara idempotent
func Run(db *gorm.DB, cfg *config.Config) {
	log.Println("===========================================")
	log.Println("Starting database seeding...")
	log.Println("===========================================")

	// Seed in correct order (respecting foreign key dependencies)
	SeedUsers(db, cfg)
	SeedProvinces(db)
	SeedRegencies(db)
	DistrictSeeder(db)
	VillageSeeder(db)
	SeedCategories(db)
	SeedBalais(db)
	RusunSeeder(db)

	log.Println("===========================================")
	log.Println("Database seeding completed successfully")
	log.Println("===========================================")
}
