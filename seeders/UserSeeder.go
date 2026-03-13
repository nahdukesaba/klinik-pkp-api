package seeders

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/user"
	"log"
)

func UserSeeder(db *gorm.DB) error {
	err := createUserIfNotExists(db, user.User{
		ID:       uuid.New(),
		Name:     "Administrator",
		Email:    "admin@klinik.com",
		Password: "admin123",
		Phone:    "081234567890",
		NIP:	  "199310202025061010",
		Role:     "admin",
		IsActive: true,
	})

	err = createUserIfNotExists(db, user.User{
		ID:       uuid.New(),
		Name:     "User Demo",
		Email:    "user@klinik.com",
		Password: "user123",
		Phone:    "081234567891",
		NIP:	  "199310202025061011",
		Role:     "user",
		IsActive: true,
	})

	return err
}

func createUserIfNotExists(db *gorm.DB, new_user user.User) error {
	var count int64

	db.Model(&user.User{}).Where("email = ?", new_user.Email).Count(&count)

	if count == 0 {
		if err := db.Create(&new_user).Error; err != nil {
			log.Printf("✗%v", err)

			return err
		} else {
			log.Printf("✓ User seeded: %s (role: %s)", new_user.Email, new_user.Role)
		}
	} else {
		log.Printf("!? User already exists: %s", new_user.Email)
	}

	return nil
}
