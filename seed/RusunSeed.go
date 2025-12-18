package seed

import (
	"klinik-api/models"
	"log"

	"gorm.io/gorm"
)

// RusunSeeder seeds rusun data for Sumatera Utara only
func RusunSeeder(db *gorm.DB) {
	log.Println("Seeding rusun...")

	// Get Sumatera Utara province and Medan regency
	var province models.Province
	var regency models.Regency
	
	if err := db.Where("code = ?", "12").First(&province).Error; err != nil {
		log.Printf("⚠ Province Sumatera Utara not found, skipping rusun seeding")
		return
	}
	
	if err := db.Where("code = ?", "1271").First(&regency).Error; err != nil {
		log.Printf("⚠ Kota Medan not found, skipping rusun seeding")
		return
	}

	// Get districts and villages
	var districtBarat, districtBelawan, districtHelvetia, districtSunggal models.District
	var districtLabuhan, districtKota, districtDeli models.District
	
	db.Where("code = ?", "1271010").First(&districtBarat)
	db.Where("code = ?", "1271020").First(&districtBelawan)
	db.Where("code = ?", "1271030").First(&districtHelvetia)
	db.Where("code = ?", "1271040").First(&districtSunggal)
	db.Where("code = ?", "1271050").First(&districtLabuhan)
	db.Where("code = ?", "1271060").First(&districtKota)
	db.Where("code = ?", "1271070").First(&districtDeli)

	var villageSeiAgul, villageSicanang, villageHelvetia, villageCintaDamai models.Village
	var villageTanjungRejo, villageMartubung, villagePasarMerahBarat, villageTanjungMulia models.Village
	
	db.Where("code = ?", "1271010001").First(&villageSeiAgul)
	db.Where("code = ?", "1271020001").First(&villageSicanang)
	db.Where("code = ?", "1271030001").First(&villageHelvetia)
	db.Where("code = ?", "1271030002").First(&villageCintaDamai)
	db.Where("code = ?", "1271040001").First(&villageTanjungRejo)
	db.Where("code = ?", "1271050001").First(&villageMartubung)
	db.Where("code = ?", "1271060001").First(&villagePasarMerahBarat)
	db.Where("code = ?", "1271070001").First(&villageTanjungMulia)

	rusunData := []models.Rusun{
		{
			Name: "Rusun Kejaksaan Tinggi", Address: "Jl. Karya Rakyat",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtBarat.ID, VillageID: villageSeiAgul.ID,
			Latitude: 3.6065, Longitude: 98.6624,
			UnitCount: 16, TowerCount: 1, FloorCount: 2, Type: "Wisama Suralaya 36",
			BuildYear: 2020, HandoverYear: 2021, OccupiedUnits: 16,
			Contractor: "PT. Pembangunan Perumahan", Status: "active",
		},
		{
			Name: "Rusun Sicanang", Address: "Jl. Sicanang No. 1, Belawan, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtBelawan.ID, VillageID: villageSicanang.ID,
			Latitude: 3.7817, Longitude: 98.6871,
			UnitCount: 96, TowerCount: 2, FloorCount: 5, Type: "Type 36",
			BuildYear: 2019, HandoverYear: 2020, OccupiedUnits: 88,
			Contractor: "PT. Adhi Karya", Status: "active",
		},
		{
			Name: "Rusun Helvetia", Address: "Jl. Karya Wisata, Helvetia, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtHelvetia.ID, VillageID: villageHelvetia.ID,
			Latitude: 3.6204, Longitude: 98.6419,
			UnitCount: 160, TowerCount: 4, FloorCount: 5, Type: "Type 24 dan 36",
			BuildYear: 2018, HandoverYear: 2019, OccupiedUnits: 152,
			Contractor: "PT. Waskita Karya", Status: "active",
		},
		{
			Name: "Rusun Tanjung Rejo", Address: "Jl. Tanjung Rejo, Medan Sunggal, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtSunggal.ID, VillageID: villageTanjungRejo.ID,
			Latitude: 3.5615, Longitude: 98.6892,
			UnitCount: 120, TowerCount: 3, FloorCount: 5, Type: "Type 36",
			BuildYear: 2020, HandoverYear: 2021, OccupiedUnits: 110,
			Contractor: "PT. Pembangunan Perumahan", Status: "active",
		},
		{
			Name: "Rusun Martubung", Address: "Jl. Veteran Pasar IV, Martubung, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtLabuhan.ID, VillageID: villageMartubung.ID,
			Latitude: 3.6558, Longitude: 98.7283,
			UnitCount: 80, TowerCount: 2, FloorCount: 4, Type: "Type 24",
			BuildYear: 2017, HandoverYear: 2018, OccupiedUnits: 76,
			Contractor: "PT. Nindya Karya", Status: "active",
		},
		{
			Name: "Rusun Cinta Damai", Address: "Jl. Cinta Damai, Medan Helvetia, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtHelvetia.ID, VillageID: villageCintaDamai.ID,
			Latitude: 3.6089, Longitude: 98.6512,
			UnitCount: 144, TowerCount: 3, FloorCount: 6, Type: "Type 36 dan 42",
			BuildYear: 2019, HandoverYear: 2020, OccupiedUnits: 138,
			Contractor: "PT. Wijaya Karya", Status: "active",
		},
		{
			Name: "Rusun Pasar Merah Barat", Address: "Jl. Pasar Merah Barat, Medan Kota, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtKota.ID, VillageID: villagePasarMerahBarat.ID,
			Latitude: 3.5778, Longitude: 98.6831,
			UnitCount: 200, TowerCount: 5, FloorCount: 5, Type: "Type 24 dan 36",
			BuildYear: 2018, HandoverYear: 2019, OccupiedUnits: 195,
			Contractor: "PT. Hutama Karya", Status: "active",
		},
		{
			Name: "Rusun Tanjung Mulia", Address: "Jl. Tanjung Mulia Hilir, Medan Deli, Medan",
			ProvinceID: province.ID, RegencyID: regency.ID,
			DistrictID: districtDeli.ID, VillageID: villageTanjungMulia.ID,
			Latitude: 3.6428, Longitude: 98.7015,
			UnitCount: 96, TowerCount: 2, FloorCount: 6, Type: "Type 36",
			BuildYear: 2020, HandoverYear: 2021, OccupiedUnits: 89,
			Contractor: "PT. Brantas Abipraya", Status: "active",
		},
	}

	for _, rusun := range rusunData {
		var count int64
		db.Model(&models.Rusun{}).Where("name = ?", rusun.Name).Count(&count)
		
		if count == 0 {
			if err := db.Create(&rusun).Error; err != nil {
				log.Printf("Failed to seed rusun %s: %v", rusun.Name, err)
			} else {
				log.Printf("✓ Rusun seeded: %s", rusun.Name)
			}
		} else {
			log.Printf("✓ Rusun already exists: %s", rusun.Name)
		}
	}

	log.Println("Rusun seeding completed")
}
