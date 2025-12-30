package seeders

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/config"
	"log"
	"strings"
	// "gorm.io/gorm/logger"
)

// RUN SEED FOR ALL MODELS OR SPECIFIC MODELS
func Run(db *gorm.DB, cfg *config.Config, seedNames ...string) {
	// UNCOMMENT TO ENABLE DETAILED LOGGING
	// db.Logger = db.Logger.LogMode(logger.Info)

	err := error(nil)

	// IF NO SEEDS SPECIFIED, SEED ALL
	if len(seedNames) == 0 {
		err = seedAll(db, cfg)
	} else {
		for _, name := range seedNames {
			err = seedOnly(db, cfg, strings.ToLower(name))

			if err != nil {
				break
			}
		}
	}

	if err == nil {
		log.Println("Database seeding completed")
	} else {
		log.Println("✗", err)
	}
}

// DISPLAY AVAILABLE SEEDERS
func DisplayAvailableSeeders() {
	fmt.Println("Available seeders:")
	fmt.Println("✓ province")
	fmt.Println("✓ region")
	fmt.Println("✓ district")
	fmt.Println("✓ village")
	fmt.Println("✓ user")
	fmt.Println("✓ rusun")

	fmt.Println("\nUsage: ")
	fmt.Println("> go run main.go seed [seed1] [seed2] ...")

	fmt.Println("\nExample: ")
	fmt.Println("> go run main.go seed province region district")
	fmt.Println("> go run main.go seed user (seeds only user table)")
	fmt.Println("> go run main.go seed (seeds all tables)")
}

// SEED ALL MODELS
func seedAll(db *gorm.DB, cfg *config.Config) error {
	err := ProvinceSeeder(db)
	err = RegionSeeder(db)
	err = DistrictSeeder(db)
	err = VillageSeeder(db)
	err = UserSeeder(db, cfg)
	err = RusunSeeder(db)

	return err
}

// SEED A SPECIFIC MODEL
func seedOnly(db *gorm.DB, cfg *config.Config, modelName string) error {
	switch modelName {
	case "province", "provinces":
		if err := ProvinceSeeder(db); err != nil {
			return err
		}
	case "region", "regions":
		if err := RegionSeeder(db); err != nil {
			return err
		}
	case "district", "districts":
		if err := DistrictSeeder(db); err != nil {
			return err
		}
	case "village", "villages":
		if err := VillageSeeder(db); err != nil {
			return err
		}
	case "user", "users":
		if err := UserSeeder(db, cfg); err != nil {
			return err
		}
	case "rusun", "rusuns":
		if err := RusunSeeder(db); err != nil {
			return err
		}
	default:
		log.Printf("✗ Unknown seed %s. Available seeders: province, region, district, village, user, rusun\n", modelName)

		return fmt.Errorf("✗ Unknown seed %s", modelName)
	}

	return nil
}
