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
	"klinik-pkp-api/internal/api/sosialisasi"
	"klinik-pkp-api/internal/api/user"
	"klinik-pkp-api/internal/api/village"
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
			&sosialisasi.Sosialisasi{},
		)

		if err != nil {
			return err
		}

		log.Println("✓ User table migrated")
		log.Println("✓ Province table migrated")
		log.Println("✓ Region table migrated")
		log.Println("✓ District table migrated")
		log.Println("✓ Village table migrated")
		log.Println("✓ Rusun table migrated")
		log.Println("✓ BSPS table migrated")
		log.Println("✓ Sosialisasi table migrated")
	} else {
		// CREATE SPECIFIC MODELS
		for _, model := range models {
			switch strings.ToLower(model) {
			case "user", "users":
				if err := db.AutoMigrate(&user.User{}); err != nil {
					return err
				}

				log.Println("✓ User table migrated")
			case "province":
				if err := db.AutoMigrate(&province.Province{}); err != nil {
					return err
				}

				log.Println("✓ Province table migrated")
			case "region":
				if err := db.AutoMigrate(&region.Region{}); err != nil {
					return err
				}

				log.Println("✓ Region table migrated")
			case "district":
				if err := db.AutoMigrate(&district.District{}); err != nil {
					return err
				}

				log.Println("✓ District table migrated")
			case "village":
				if err := db.AutoMigrate(&village.Village{}); err != nil {
					return err
				}

				log.Println("✓ Village table migrated")
			case "rusun":
				if err := db.AutoMigrate(&rusun.Rusun{}); err != nil {
					return err
				}

				log.Println("✓ Rusun table migrated")
			case "bsps":
				if err := db.AutoMigrate(&bsps.BSPS{}); err != nil {
					return err
				}

				log.Println("✓ BSPS table migrated")
			case "sosialisasi":
				if err := db.AutoMigrate(&sosialisasi.Sosialisasi{}); err != nil {
					return err
				}

				log.Println("✓ Sosialisasi table migrated")
			default:
				return fmt.Errorf("Unknown migration %s. Available migrations: user, rusun, bsps, province, region, district, village\n", strings.ToLower(model))
			}
		}
	}

	log.Println("Table(s) creation completed")

	return err
}

// DROP TABLES AND RESET SEQUENCES BASED ON TABLES ARGUMENTS
func dropTables(db *gorm.DB, tables ...string) error {
	var err error

	// IF NO TABLES SPECIFIED, DROP ALL
	if len(tables) == 0 {
		var tablesFromDatabase []string

		if err = db.Raw(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`).Scan(&tablesFromDatabase).Error; err != nil {
			return fmt.Errorf("failed to get table names: %v", err)
		}

		if len(tablesFromDatabase) == 0 {
			return fmt.Errorf("no tables found in database")
		}

		for _, table := range tablesFromDatabase {
			if err = utils.DropTable(db, table); err != nil {
				return err
			}
		}
	} else {
		// DROP SPECIFIC TABLES
		for _, table := range tables {
			if err = utils.DropTable(db, table); err != nil {
				return err
			}
		}
	}

	log.Println("Table(s) deletion completed")

	return nil
}

// TRUNCATE A SPECIFIC TABLE BASED ON TABLES ARGUMENT
func truncateTables(db *gorm.DB, tables ...string) error {
	var err error

	// IF NO TABLES SPECIFIED, TRUNCATE ALL
	if len(tables) == 0 {
		var tablesFromDatabase []string

		if err = db.Raw(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`).Scan(&tablesFromDatabase).Error; err != nil {
			return fmt.Errorf("%v", err)
		}

		if len(tablesFromDatabase) == 0 {
			return fmt.Errorf("no tables found in database")
		}

		for _, table := range tablesFromDatabase {
			if err = utils.TruncateTable(db, table); err != nil {
				return err
			}
		}
	} else {
		// TRUNCATE SPECIFIC TABLES
		for _, table := range tables {
			if err = utils.TruncateTable(db, table); err != nil {
				return err
			}
		}
	}

	log.Printf("Table(s) truncation completed")

	return nil
}

// RUN MODIFICATION BASED ON MIGRATION FILE
func modifyTable(db *gorm.DB, filename string) error {
	if err := utils.RunModification(db, filename); err != nil {
		return fmt.Errorf("failed to run migration %s: %v", filename, err)
	}

	log.Printf("Table(s) modification completed")

	return nil
}

// ROLLBACK A SPECIFIC MODIFICATION BASED ON MIGRATION FILE
func rollbackTable(db *gorm.DB, filename string) error {
	if err := utils.RunRollback(db, filename); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %v", filename, err)
	}

	log.Printf("Table(s) rollback completed")

	return nil
}

// DISPLAY AVAILABLE MIGRATIONS
func displayAvailableMigrations() {
	fmt.Println("1. Creating table(s):")
	fmt.Println("   > go run cmd/migrate/main.go")
	fmt.Println("   > go run cmd/migrate/main.go user rusun bsps province region district village")
	fmt.Println()

	fmt.Println("2. Dropping table(s):")
	fmt.Println("   > go run cmd/migrate/main.go drop")
	fmt.Println("   > go run cmd/migrate/main.go drop user rusun bsps")
	fmt.Println()

	fmt.Println("3. Truncating table(s):")
	fmt.Println("   > go run cmd/migrate/main.go truncate")
	fmt.Println("   > go run cmd/migrate/main.go truncate user rusun bsps")
	fmt.Println()

	fmt.Println("4. Modifying table(s):")
	fmt.Println("   > go run cmd/migrate/main.go modify 20260104_add_example_column_to_bsps")
	fmt.Println("   > go run cmd/migrate/main.go modify 20260104_add_indexes_to_bsps")
	fmt.Println()

	fmt.Println("5. Rolling back table(s):")
	fmt.Println("   > go run cmd/migrate/main.go rollback 20260104_add_example_column_to_bsps")
	fmt.Println()
}

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	command := ""

	if len(os.Args) >= 2 {
		command = strings.ToLower(os.Args[1])
	}

	// HANDLE EXECUTION BASED ON FLAGS
	if len(os.Args) == 2 && (command == "--help" || command == "-h" || command == "help") {
		displayAvailableMigrations()

		return
	} else if len(os.Args) >= 2 && (command == "drop" || command == "truncate" || command == "rollback" || command == "modify") {
		switch command {
		case "drop":
			if err := dropTables(db, os.Args[2:]...); err != nil {
				log.Fatal("✗ Failed to drop tables: ", err)
			}
		case "truncate":
			if err := truncateTables(db, os.Args[2:]...); err != nil {
				log.Fatal("✗ Failed to truncate table: ", err)
			}
		case "rollback":
			if len(os.Args) < 3 {
				log.Fatal("✗ Please specify a file name. Usage: go run cmd/migrate/main.go rollback <file_name>")
			}

			if err := rollbackTable(db, os.Args[2]); err != nil {
				log.Fatal("✗ Failed to rollback migration: ", err)
			}
		case "modify":
			if len(os.Args) < 3 {
				log.Fatal("✗ Please specify a file name. Usage: go run cmd/migrate/main.go modify <file_name>")
			}

			if err := modifyTable(db, os.Args[2]); err != nil {
				log.Fatal("✗ Failed to run migration: ", err)
			}
		}
	} else {
		if err := createTables(db, os.Args[1:]...); err != nil {
			log.Println("✗", err)
		}
	}
}
