package main

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/bsps"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/province"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/rusun"
	"klinik-pkp-api/internal/api/user"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/migrations"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"strings"
)

// CREATE TABLES BASED ON MODELS ARGUMENTS
func createTables(db *gorm.DB, models ...string) error {
	var err error

	// IF NO MODELS SPECIFIED, CREATE ALL
	if len(models) == 0 {
		err = db.AutoMigrate(
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
	} else {
		// CREATE SPECIFIC MODELS
		for _, model := range models {
			switch strings.ToLower(model) {
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
				return fmt.Errorf("Unknown migration %s. Available migrations: user, rusun, bsps, province, region, district, village\n", strings.ToLower(model))
			}
		}
	}

	return err
}

// DROP TABLES AND RESET SEQUENCES BASED ON MODELS ARGUMENTS
func dropTables(db *gorm.DB, tables ...string) error {
	var err error

	// IF NO TABLES SPECIFIED, DROP ALL
	if len(tables) == 0 {
		var tablesFromDatabase []string

		if err = db.Raw(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`).Scan(&tablesFromDatabase).Error; err != nil {
			return fmt.Errorf("failed to get table names: %v", err)
		}

		for _, table := range tablesFromDatabase {
			err = utils.DropTable(db, table)
		}
	} else {
		// DROP SPECIFIC TABLES
		for _, table := range tables {
			err = utils.DropTable(db, table)
		}
	}

	log.Println("Table(s) deletion completed")

	return err
}

// TRUNCATE A SPECIFIC TABLE BASED ON MODELS ARGUMENT
func truncateTables(db *gorm.DB, tables ...string) error {
	var err error

	// IF NO TABLES SPECIFIED, TRUNCATE ALL
	if len(tables) == 0 {
		var tablesFromDatabase []string

		if err = db.Raw(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`).Scan(&tablesFromDatabase).Error; err != nil {
			return fmt.Errorf("failed to get table names: %v", err)
		}

		for _, table := range tablesFromDatabase {
			err = utils.TruncateTable(db, table)
		}
	} else {
		// DROP SPECIFIC TABLES
		for _, table := range tables {
			err = utils.TruncateTable(db, table)
		}
	}

	log.Printf("✓ Table(s) truncated successfully")

	return err
}

// RUN MODIFICATION BASED ON MIGRATION FILE
func modifyTable(db *gorm.DB, filename string) error {
	log.Printf("Looking for migration: %s", filename)

	if err := migrations.RunModification(db, filename); err != nil {
		return fmt.Errorf("failed to run migration %s: %v", filename, err)
	}

	log.Printf("✓ Migration %s completed successfully", filename)

	return nil
}

// ROLLBACK A SPECIFIC MODIFICATION BASED ON MIGRATION FILE
func rollbackTable(db *gorm.DB, filename string) error {
	log.Printf("Looking for rollback: %s", filename)

	if err := migrations.RunRollback(db, filename); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %v", filename, err)
	}

	log.Printf("✓ Rollback %s completed successfully", filename)

	return nil
}

// DISPLAY AVAILABLE MIGRATIONS
func displayAvailableMigrations() {
	fmt.Println("1. Migrate all tables:")
	fmt.Println("   > go run cmd/migrate/main.go")
	fmt.Println()

	fmt.Println("2. Migrate specific tables:")
	fmt.Println("   > go run cmd/migrate/main.go user province region")
	fmt.Println()

	fmt.Println("3. Drop all tables and sequences:")
	fmt.Println("   > go run cmd/migrate/main.go drop")
	fmt.Println()

	fmt.Println("4. Truncate a specific table:")
	fmt.Println("   > go run cmd/migrate/main.go truncate bsps")
	fmt.Println("   > go run cmd/migrate/main.go truncate user")
	fmt.Println()

	fmt.Println("5. Run custom migration file:")
	fmt.Println("   > go run cmd/migrate/main.go 20260104_add_example_column_to_bsps")
	fmt.Println("   > go run cmd/migrate/main.go 20260104_add_indexes_to_bsps")
	fmt.Println()

	fmt.Println("6. Rollback a migration:")
	fmt.Println("   > go run cmd/migrate/main.go rollback 20260104_add_example_column_to_bsps")
	fmt.Println()

	fmt.Println("📦 Available Models:")
	fmt.Println("   ✓ user      - User accounts")
	fmt.Println("   ✓ rusun     - Rusun data")
	fmt.Println("   ✓ province  - Province data")
	fmt.Println("   ✓ region    - Region data")
	fmt.Println("   ✓ district  - District data")
	fmt.Println("   ✓ village   - Village data")
	fmt.Println("   ✓ bsps      - BSPS data")
	fmt.Println()
}

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	// HANDLE EXECUTION BASED ON FLAGS
	if len(os.Args) == 2 && (strings.ToLower(os.Args[1]) == "--help" || strings.ToLower(os.Args[1]) == "-h" || strings.ToLower(os.Args[1]) == "help") {
		displayAvailableMigrations()

		return
	} else if len(os.Args) >= 2 && (strings.ToLower(os.Args[1]) == "drop" || strings.ToLower(os.Args[1]) == "truncate" || strings.ToLower(os.Args[1]) == "rollback" || strings.ToLower(os.Args[1]) == "modify") {
		switch strings.ToLower(os.Args[1]) {
		case "drop":
			if err := dropTables(db, os.Args[2:]...); err != nil {
				log.Fatal("Failed to drop tables:", err)
			}
		case "truncate":
			if err := truncateTables(db, os.Args[2:]...); err != nil {
				log.Fatal("Failed to truncate table:", err)
			}
		case "rollback":
			if len(os.Args) < 3 {
				log.Fatal("Please specify a file name. Usage: go run cmd/migrate/main.go rollback <file_name>")
			}

			fileName := os.Args[2]

			if err := rollbackTable(db, fileName); err != nil {
				log.Fatal("Failed to rollback migration:", err)
			}
		case "modify":
			if len(os.Args) < 3 {
				log.Fatal("Please specify a file name. Usage: go run cmd/migrate/main.go modify <file_name>")
			}

			fileName := os.Args[2]

			if err := modifyTable(db, fileName); err != nil {
				log.Fatal("Failed to run migration:", err)
			}
		}
	} else {
		if err := createTables(db, os.Args[1:]...); err != nil {
			log.Println("✗", err)
		} else {
			log.Println("Database creation completed")
		}
	}
}
