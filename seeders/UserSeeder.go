package seeders

import (
	"klinik-pkp-api/config"
	"klinik-pkp-api/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SEEDS FROM PREDEFINED DATA
func UserSeeder(db *gorm.DB, cfg *config.Config) error {
	err := createUserIfNotExists(db, models.User{
		ID:       uuid.New(),
		Name:     "Administrator",
		Email:    "admin@klinik.com",
		Password: "admin123",
		Phone:    "081234567890",
		Role:     "admin",
		IsActive: true,
	})

	err = createUserIfNotExists(db, models.User{
		ID:       uuid.New(),
		Name:     "User Demo",
		Email:    "user@klinik.com",
		Password: "user123",
		Phone:    "081234567891",
		Role:     "user",
		IsActive: true,
	})

	return err
}

func createUserIfNotExists(db *gorm.DB, user models.User) error {
	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)

	if count == 0 {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("✗%v", err)
			return err
		} else {
			log.Printf("✓ User seeded: %s (role: %s)", user.Email, user.Role)
		}
	} else {
		log.Printf("User already exists: %s", user.Email)
	}

	return nil
}
