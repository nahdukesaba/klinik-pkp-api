package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// VillageSeeder seeds village (kelurahan) data for Medan districts
func VillageSeeder(db *gorm.DB) {
	log.Println("Seeding villages...")

	// Get districts
	var districtBarat, districtBelawan, districtHelvetia, districtSunggal models.District
	var districtLabuhan, districtKota, districtDeli models.District
	
	db.Where("code = ?", "1271010").First(&districtBarat)
	db.Where("code = ?", "1271020").First(&districtBelawan)
	db.Where("code = ?", "1271030").First(&districtHelvetia)
	db.Where("code = ?", "1271040").First(&districtSunggal)
	db.Where("code = ?", "1271050").First(&districtLabuhan)
	db.Where("code = ?", "1271060").First(&districtKota)
	db.Where("code = ?", "1271070").First(&districtDeli)

	villages := []models.Village{
		// Medan Barat
		{DistrictID: districtBarat.ID, Name: "Sei Agul", Code: "1271010001"},
		
		// Medan Belawan
		{DistrictID: districtBelawan.ID, Name: "Sicanang", Code: "1271020001"},
		
		// Medan Helvetia
		{DistrictID: districtHelvetia.ID, Name: "Helvetia", Code: "1271030001"},
		{DistrictID: districtHelvetia.ID, Name: "Cinta Damai", Code: "1271030002"},
		
		// Medan Sunggal
		{DistrictID: districtSunggal.ID, Name: "Tanjung Rejo", Code: "1271040001"},
		
		// Medan Labuhan
		{DistrictID: districtLabuhan.ID, Name: "Martubung", Code: "1271050001"},
		
		// Medan Kota
		{DistrictID: districtKota.ID, Name: "Pasar Merah Barat", Code: "1271060001"},
		
		// Medan Deli
		{DistrictID: districtDeli.ID, Name: "Tanjung Mulia Hilir", Code: "1271070001"},
	}

	for _, village := range villages {
		var count int64
		db.Model(&models.Village{}).Where("code = ?", village.Code).Count(&count)
		
		if count == 0 {
			if err := db.Create(&village).Error; err != nil {
				log.Printf("Failed to seed village %s: %v", village.Name, err)
			} else {
				log.Printf("✓ Village seeded: %s", village.Name)
			}
		} else {
			log.Printf("✓ Village already exists: %s", village.Name)
		}
	}

	log.Println("Villages seeding completed")
}
