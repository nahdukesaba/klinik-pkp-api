package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	// "gorm.io/gorm/logger"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	
	//	[
	//  	{"latitude":14.23,"longitude":120.98},
	//   	{"latitude":14.25,"longitude":121.01}
	//	]
}

// RUN MIGRATION FOR ALL MODELS OR SPECIFIC MODEL
func Run(db *gorm.DB, models ...string) {
	// UNCOMMENT TO ENABLE DETAILED LOGGING
	// db.Logger = db.Logger.LogMode(logger.Info)

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
func DisplayAvailableMigrations() {
	fmt.Println("Available migrations:")
	fmt.Println("✓ user")
	fmt.Println("✓ rusun")
	fmt.Println("✓ province")
	fmt.Println("✓ region")
	fmt.Println("✓ district")
	fmt.Println("✓ village")

	fmt.Println("\nUsage: ")
	fmt.Println("> go run main.go migrate [model1] [model2] ...")

	fmt.Println("\nExample: ")
	fmt.Println("> go run main.go migrate province region district")
	fmt.Println("> go run main.go migrate user (migrates only user table)")
	fmt.Println("> go run main.go migrate (migrates all tables)")
}

// MIGRATE ALL MODELS
func migrateAll(db *gorm.DB) error {
	err := db.AutoMigrate(
		&Province{},
		&Region{},
		&District{},
		&Village{},
		&User{},
		&Rusun{},
		&BSPS{},
	)

	if err != nil {
		return err
	}

	log.Println("✓ User table migrated")
	log.Println("✓ Rusun table migrated")
	log.Println("✓ BSPS table migrated")
	log.Println("✓ Province table migrated")
	log.Println("✓ Region table migrated")
	log.Println("✓ District table migrated")
	log.Println("✓ Village table migrated")

	return err
}

// MIGRATE A SPECIFIC MODEL
func migrateOnly(db *gorm.DB, modelName string) error {
	switch modelName {
	case "user", "users":
		if err := db.AutoMigrate(&User{}); err != nil {
			return err
		}

		log.Println("✓ User table migrated")
	case "rusun", "rusuns":
		if err := db.AutoMigrate(&Rusun{}); err != nil {
			return err
		}

		log.Println("✓ Rusun table migrated")
	case "bsps":
		if err := db.AutoMigrate(&BSPS{}); err != nil {
			return err
		}

		log.Println("✓ BSPS table migrated")
	case "province", "provinces":
		if err := db.AutoMigrate(&Province{}); err != nil {
			return err
		}

		log.Println("✓ Province table migrated")
	case "region", "regions":
		if err := db.AutoMigrate(&Region{}); err != nil {
			return err
		}

		log.Println("✓ Region table migrated")
	case "district", "districts":
		if err := db.AutoMigrate(&District{}); err != nil {
			return err
		}

		log.Println("✓ District table migrated")
	case "village", "villages":
		if err := db.AutoMigrate(&Village{}); err != nil {
			return err
		}

		log.Println("✓ Village table migrated")
	default:
		log.Printf("✗ Unknown migration %s. Available migrations: user, rusun, province, region, district, village\n", modelName)

		return fmt.Errorf("✗ Unknown migration %s", modelName)
	}

	return nil
}
