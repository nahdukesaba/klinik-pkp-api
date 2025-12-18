package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// SeedBalais seeds balai data
func SeedBalais(db *gorm.DB) {
	log.Println("Seeding balais...")

	// Get required IDs
	var sumutProvince models.Province
	var jalanCategory models.Category

	db.Where("code = ?", "12").First(&sumutProvince)
	db.Where("name = ?", "Jalan").First(&jalanCategory)

	if sumutProvince.ID == 0 || jalanCategory.ID == 0 {
		log.Println("⚠ Required province or category not found, skipping balais seeding")
		return
	}

	balais := []models.Balai{
		{
			Name:       "Balai Besar Pelaksanaan Jalan Nasional I",
			Address:    "Jl. Sisingamangaraja No.1, Medan",
			ProvinceID: sumutProvince.ID,
			CategoryID: jalanCategory.ID,
		},
		{
			Name:       "Balai Besar Pelaksanaan Jalan Nasional II",
			Address:    "Jl. Gatot Subroto No.45, Medan",
			ProvinceID: sumutProvince.ID,
			CategoryID: jalanCategory.ID,
		},
	}

	for _, balai := range balais {
		var count int64
		db.Model(&models.Balai{}).Where("name = ?", balai.Name).Count(&count)

		if count == 0 {
			if err := db.Create(&balai).Error; err != nil {
				log.Printf("Failed to seed balai %s: %v", balai.Name, err)
			} else {
				log.Printf("✓ Balai seeded: %s", balai.Name)
			}
		} else {
			log.Printf("✓ Balai already exists: %s", balai.Name)
		}
	}

	log.Println("Balais seeding completed")
}
