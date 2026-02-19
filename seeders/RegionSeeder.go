package seeders

import (
	"encoding/csv"
	"klinik-pkp-api/internal/api/region"
	"log"
	"os"
	"strconv"
	"gorm.io/gorm"
)

// SEEDS FROM THE CSV FILE
func RegionSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/region.csv")

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
	var data []region.Region

	db.Model(&region.Region{}).Count(&count)

	if count > 0 {
		log.Println("!? Region table already seeded, skipping...")

		return nil
	}

	// _, records[1:] SKIP HEADER ROW
	for _, record := range records[1:] {
		if len(record) < 3 {
			continue
		}

		id, err := strconv.ParseUint(record[0], 10, 64)

		if err != nil {
			log.Println("Error parsing ID:", err)

			continue
		}

		provinceID, err := strconv.ParseUint(record[2], 10, 64)

		if err != nil {
			log.Println("Error parsing ProvinceID:", err)

			continue
		}

		region := region.Region{
			ID:         id,
			Name:       record[1],
			ProvinceID: provinceID,
		}

		data = append(data, region)
	}

	if len(data) > 0 {
		result := db.CreateInBatches(data, 100)

		if result.Error != nil {
			return result.Error
		}

		log.Println("✓ Region table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
