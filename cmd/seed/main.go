package main

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/config"
	"klinik-pkp-api/seeders"
	"log"
	"os"
	"strings"
)

// RUN SEED FOR ALL MODELS OR SPECIFIC MODELS
func runSeeders(db *gorm.DB, seedNames ...string) {
	err := error(nil)

	// IF NO SEEDS SPECIFIED, SEED ALL
	if len(seedNames) == 0 {
		err = seedAll(db)
	} else {
		for _, name := range seedNames {
			err = seedOnly(db, strings.ToLower(name))

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
func displayAvailableSeeders() {
	fmt.Println("Available seeders:")
	fmt.Println("✓ province")
	fmt.Println("✓ region")
	fmt.Println("✓ district")
	fmt.Println("✓ village")
	fmt.Println("✓ user")
	fmt.Println("✓ rusun")
	fmt.Println("✓ bsps")

	fmt.Println("\nUsage: ")
	fmt.Println("> go run cmd/seed/main.go [seed1] [seed2] ...")

	fmt.Println("\nExample: ")
	fmt.Println("> go run cmd/seed/main.go province region district")
	fmt.Println("> go run cmd/seed/main.go user (seeds only user table)")
	fmt.Println("> go run cmd/seed/main.go (seeds all tables)")
}

// SEED ALL MODELS
func seedAll(db *gorm.DB) error {
	err := seeders.ProvinceSeeder(db)
	err = seeders.RegionSeeder(db)
	err = seeders.DistrictSeeder(db)
	err = seeders.VillageSeeder(db)
	// err = seeders.UserSeeder(db)
	err = seeders.RusunSeeder(db)
	err = seeders.BSPSSeeder(db)

	return err
}

// SEED A SPECIFIC MODEL
func seedOnly(db *gorm.DB, modelName string) error {
	switch modelName {
	case "province", "provinces":
		if err := seeders.ProvinceSeeder(db); err != nil {
			return err
		}
	case "region", "regions":
		if err := seeders.RegionSeeder(db); err != nil {
			return err
		}
	case "district", "districts":
		if err := seeders.DistrictSeeder(db); err != nil {
			return err
		}
	case "village", "villages":
		if err := seeders.VillageSeeder(db); err != nil {
			return err
		}
	case "user", "users":
		if err := seeders.UserSeeder(db); err != nil {
			return err
		}
	case "rusun", "rusuns":
		if err := seeders.RusunSeeder(db); err != nil {
			return err
		}
	case "bsps":
		if err := seeders.BSPSSeeder(db); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown seed %s. Available seeders: province, region, district, village, user, rusun, bsps\n", modelName)
	}

	return nil
}

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	// HANDLE HELP FLAG
	if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		displayAvailableSeeders()

		return
	} else {
		// RUN SEEDER
		runSeeders(db, os.Args[1:]...)
	}
}
