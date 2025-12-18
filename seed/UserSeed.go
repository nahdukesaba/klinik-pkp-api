package seed

import (
	"klinik-api/config"
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// SeedUsers seeds all user data with different roles
func SeedUsers(db *gorm.DB, cfg *config.Config) {
	log.Println("Seeding users...")

	// Admin user
	seedUser(db, models.User{
		Name:     "Administrator",
		Email:    "admin@klinik.com",
		Password: "admin123",
		Phone:    "081234567890",
		Role:     "admin",
		IsActive: true,
	})

	// Regular user
	seedUser(db, models.User{
		Name:     "User Demo",
		Email:    "user@klinik.com",
		Password: "user123",
		Phone:    "081234567891",
		Role:     "user",
		IsActive: true,
	})

	log.Println("Users seeding completed")
}

func seedUser(db *gorm.DB, user models.User) {
	var count int64
	db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)

	if count == 0 {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Email, err)
		} else {
			log.Printf("✓ User seeded: %s (role: %s)", user.Email, user.Role)
		}
	} else {
		log.Printf("✓ User already exists: %s", user.Email)
	}
}
