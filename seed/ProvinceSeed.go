package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// SeedProvinces seeds province data
func SeedProvinces(db *gorm.DB) {
	log.Println("Seeding provinces...")

	provinces := []models.Province{
		{Code: "12", Name: "Sumatera Utara"},
		{Code: "11", Name: "Aceh"},
		{Code: "13", Name: "Sumatera Barat"},
		{Code: "14", Name: "Riau"},
		{Code: "15", Name: "Jambi"},
	}

	for _, province := range provinces {
		var count int64
		db.Model(&models.Province{}).Where("code = ?", province.Code).Count(&count)

		if count == 0 {
			if err := db.Create(&province).Error; err != nil {
				log.Printf("Failed to seed province %s: %v", province.Name, err)
			} else {
				log.Printf("✓ Province seeded: %s (code: %s)", province.Name, province.Code)
			}
		} else {
			log.Printf("✓ Province already exists: %s", province.Name)
		}
	}

	log.Println("Provinces seeding completed")
}
