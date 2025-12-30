package seeders

import (
	"encoding/csv"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/models"
	"log"
	"os"
)

// SEEDS FROM THE CSV FILE
func DistrictSeeder(db *gorm.DB) error {
	file, err := os.Open("script/district.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&models.District{}).Count(&count)

	if count > 0 {
		log.Println("!? District table already seeded, skipping...")

		return nil
	}

	reader := csv.NewReader(file)
	const batchSize = 100
	var batch []models.District

	// _ MEANS EXCLUDE THE RESULT FROM THE FIRST READ (REMOVING HEADER ROW)
	if _, err := reader.Read(); err != nil {
		return err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading district.csv:", err)
			continue
		}

		if len(record) < 3 {
			continue
		}

		batch = append(batch, models.District{
			ID:       record[0],
			Name:     record[1],
			RegionID: record[2],
		})
	}

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		batch = batch[:0] // RESET BATCH

		log.Println("✓ District table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
