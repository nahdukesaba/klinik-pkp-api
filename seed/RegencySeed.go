package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// SeedRegencies seeds regency data
func SeedRegencies(db *gorm.DB) {
	log.Println("Seeding regencies...")

	// Get province IDs
	var sumutProvince models.Province
	db.Where("code = ?", "12").First(&sumutProvince)

	if sumutProvince.ID == 0 {
		log.Println("⚠ Province 'Sumatera Utara' not found, skipping regencies seeding")
		return
	}

	regencies := []models.Regency{
		{ProvinceID: sumutProvince.ID, Code: "1201", Name: "Kabupaten Nias"},
		{ProvinceID: sumutProvince.ID, Code: "1202", Name: "Kabupaten Mandailing Natal"},
		{ProvinceID: sumutProvince.ID, Code: "1203", Name: "Kabupaten Tapanuli Selatan"},
		{ProvinceID: sumutProvince.ID, Code: "1204", Name: "Kabupaten Tapanuli Tengah"},
		{ProvinceID: sumutProvince.ID, Code: "1205", Name: "Kabupaten Tapanuli Utara"},
		{ProvinceID: sumutProvince.ID, Code: "1271", Name: "Kota Medan"},
		{ProvinceID: sumutProvince.ID, Code: "1272", Name: "Kota Binjai"},
		{ProvinceID: sumutProvince.ID, Code: "1273", Name: "Kota Pematang Siantar"},
	}

	for _, regency := range regencies {
		var count int64
		db.Model(&models.Regency{}).Where("code = ?", regency.Code).Count(&count)

		if count == 0 {
			if err := db.Create(&regency).Error; err != nil {
				log.Printf("Failed to seed regency %s: %v", regency.Name, err)
			} else {
				log.Printf("✓ Regency seeded: %s (code: %s)", regency.Name, regency.Code)
			}
		} else {
			log.Printf("✓ Regency already exists: %s", regency.Name)
		}
	}

	log.Println("Regencies seeding completed")
}
