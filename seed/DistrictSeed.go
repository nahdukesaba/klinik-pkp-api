package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// DistrictSeeder seeds district (kecamatan) data for Medan
func DistrictSeeder(db *gorm.DB) {
	log.Println("Seeding districts...")

	// Get Kota Medan regency
	var regency models.Regency
	if err := db.Where("code = ?", "1271").First(&regency).Error; err != nil {
		log.Printf("⚠ Kota Medan not found, skipping district seeding")
		return
	}

	districts := []models.District{
		{RegencyID: regency.ID, Name: "Medan Barat", Code: "1271010"},
		{RegencyID: regency.ID, Name: "Medan Belawan", Code: "1271020"},
		{RegencyID: regency.ID, Name: "Medan Helvetia", Code: "1271030"},
		{RegencyID: regency.ID, Name: "Medan Sunggal", Code: "1271040"},
		{RegencyID: regency.ID, Name: "Medan Labuhan", Code: "1271050"},
		{RegencyID: regency.ID, Name: "Medan Kota", Code: "1271060"},
		{RegencyID: regency.ID, Name: "Medan Deli", Code: "1271070"},
	}

	for _, district := range districts {
		var count int64
		db.Model(&models.District{}).Where("code = ?", district.Code).Count(&count)
		
		if count == 0 {
			if err := db.Create(&district).Error; err != nil {
				log.Printf("Failed to seed district %s: %v", district.Name, err)
			} else {
				log.Printf("✓ District seeded: %s", district.Name)
			}
		} else {
			log.Printf("✓ District already exists: %s", district.Name)
		}
	}

	log.Println("Districts seeding completed")
}
