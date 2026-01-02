package main

import (
	"fmt"
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/bsps"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/province"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/rusun"
	"klinik-pkp-api/internal/api/user"
	"klinik-pkp-api/internal/api/village"
	"log"
	"os"
	"strings"
	"gorm.io/gorm"
)

// RUN MIGRATION FOR ALL MODELS OR SPECIFIC MODEL
func runMigrations(db *gorm.DB, models ...string) {
	err := error(nil)

	// IF NO MODELS SPECIFIED, MIGRATE ALL
	if len(models) == 0 {
		err = migrateAll(db)
	} else {
		for _, model := range models {
			err = migrateOnly(db, strings.ToLower(model))
		}
	}

	if err == nil {
		log.Println("Database migration completed")
	} else {
		log.Println("✗", err)
	}
}

// DISPLAY AVAILABLE MODELS FOR MIGRATION
func displayAvailableMigrations() {
	fmt.Println("Available migrations:")
	fmt.Println("✓ user")
	fmt.Println("✓ rusun")
	fmt.Println("✓ province")
	fmt.Println("✓ region")
	fmt.Println("✓ district")
	fmt.Println("✓ village")
	fmt.Println("✓ bsps")

	fmt.Println("\nUsage: ")
	fmt.Println("> go run cmd/migrate/main.go [model1] [model2] ...")

	fmt.Println("\nExample: ")
	fmt.Println("> go run cmd/migrate/main.go province region district")
	fmt.Println("> go run cmd/migrate/main.go user (migrates only user table)")
	fmt.Println("> go run cmd/migrate/main.go (migrates all tables)")
}

// MIGRATE ALL MODELS
func migrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&province.Province{},
		&region.Region{},
		&district.District{},
		&village.Village{},
		&user.User{},
		&rusun.Rusun{},
		&bsps.BSPS{},
	)

	if err != nil {
		return err
	}

	log.Println("✓ Province table migrated")
	log.Println("✓ Region table migrated")
	log.Println("✓ District table migrated")
	log.Println("✓ Village table migrated")
	log.Println("✓ User table migrated")
	log.Println("✓ Rusun table migrated")
	log.Println("✓ BSPS table migrated")

	return err
}

// MIGRATE A SPECIFIC MODEL
func migrateOnly(db *gorm.DB, modelName string) error {
	switch modelName {
	case "user", "users":
		if err := db.AutoMigrate(&user.User{}); err != nil {
			return err
		}
		log.Println("✓ User table migrated")
	case "rusun", "rusuns":
		if err := db.AutoMigrate(&rusun.Rusun{}); err != nil {
			return err
		}
		log.Println("✓ Rusun table migrated")
	case "bsps":
		if err := db.AutoMigrate(&bsps.BSPS{}); err != nil {
			return err
		}
		log.Println("✓ BSPS table migrated")
	case "province", "provinces":
		if err := db.AutoMigrate(&province.Province{}); err != nil {
			return err
		}
		log.Println("✓ Province table migrated")
	case "region", "regions":
		if err := db.AutoMigrate(&region.Region{}); err != nil {
			return err
		}
		log.Println("✓ Region table migrated")
	case "district", "districts":
		if err := db.AutoMigrate(&district.District{}); err != nil {
			return err
		}
		log.Println("✓ District table migrated")
	case "village", "villages":
		if err := db.AutoMigrate(&village.Village{}); err != nil {
			return err
		}
		log.Println("✓ Village table migrated")
	default:
		log.Printf("✗ Unknown migration %s. Available migrations: user, rusun, bsps, province, region, district, village\n", modelName)
		return fmt.Errorf("✗ Unknown migration %s", modelName)
	}

	return nil
}

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	// HANDLE HELP FLAG
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		displayAvailableMigrations()
		return
	}

	// RUN MIGRATION
	runMigrations(db, os.Args[1:]...)
}
