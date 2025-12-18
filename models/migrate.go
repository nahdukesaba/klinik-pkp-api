package models

import (
	"log"

	"gorm.io/gorm"
)

// RunMigrations menjalankan auto migrate untuk semua models
func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&Province{},
		&Regency{},
		&District{},
		&Village{},
		&Category{},
		&Balai{},
		&Rusun{},
		&Image{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Database migrations completed successfully")
}

