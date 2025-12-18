package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// SeedCategories seeds category data
func SeedCategories(db *gorm.DB) {
	log.Println("Seeding categories...")

	categories := []models.Category{
		{
			Name:        "Jalan",
			Description: "Kategori untuk infrastruktur jalan",
		},
		{
			Name:        "Jembatan",
			Description: "Kategori untuk infrastruktur jembatan",
		},
		{
			Name:        "Gedung",
			Description: "Kategori untuk bangunan gedung",
		},
		{
			Name:        "Irigasi",
			Description: "Kategori untuk sistem irigasi",
		},
		{
			Name:        "Air Bersih",
			Description: "Kategori untuk infrastruktur air bersih",
		},
	}

	for _, category := range categories {
		var count int64
		db.Model(&models.Category{}).Where("name = ?", category.Name).Count(&count)

		if count == 0 {
			if err := db.Create(&category).Error; err != nil {
				log.Printf("Failed to seed category %s: %v", category.Name, err)
			} else {
				log.Printf("✓ Category seeded: %s", category.Name)
			}
		} else {
			log.Printf("✓ Category already exists: %s", category.Name)
		}
	}

	log.Println("Categories seeding completed")
}
