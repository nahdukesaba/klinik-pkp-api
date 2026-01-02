package seeders

import (
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/user"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SEEDS FROM PREDEFINED DATA
func UserSeeder(db *gorm.DB, cfg *config.Config) error {
	err := createUserIfNotExists(db, user.User{
		ID:       uuid.New(),
		Name:     "Administrator",
		Email:    "admin@klinik.com",
		Password: "admin123",
		Phone:    "081234567890",
		Role:     "admin",
		IsActive: true,
	})

	err = createUserIfNotExists(db, user.User{
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

func createUserIfNotExists(db *gorm.DB, u user.User) error {
	var count int64
	db.Model(&user.User{}).Where("email = ?", u.Email).Count(&count)

	if count == 0 {
		if err := db.Create(&u).Error; err != nil {
			log.Printf("✗%v", err)
			return err
		} else {
			log.Printf("✓ User seeded: %s (role: %s)", u.Email, u.Role)
		}
	} else {
		log.Printf("User already exists: %s", u.Email)
	}

	return nil
}
