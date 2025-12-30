package seeders

import (
	"encoding/csv"
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"log"
	"os"
)

// SEEDS FROM THE CSV FILE
func ProvinceSeeder(db *gorm.DB) error {
	file, err := os.Open("script/province.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return err
	}

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64
	var data []models.Province

	db.Model(&models.Province{}).Count(&count)

	if count > 0 {
		log.Println("!? Province table already seeded, skipping...")

		return nil
	}

	// _, records[1:] SKIP HEADER ROW
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		id := record[0]

		province := models.Province{
			ID:   id,
			Name: record[1],
		}

		data = append(data, province)
	}

	if len(data) > 0 {
		result := db.CreateInBatches(data, 100)

		if result.Error != nil {
			return result.Error
		}

		log.Println("✓ Province table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
