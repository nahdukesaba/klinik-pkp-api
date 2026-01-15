package seeders

import (
	"encoding/csv"
	"io"
	"klinik-pkp-api/internal/api/village"
	"log"
	"os"
	"gorm.io/gorm"
)

// SEEDS FROM THE CSV FILE
func VillageSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/village.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&village.Village{}).Count(&count)

	if count > 0 {
		log.Println("!? Village table already seeded, skipping...")

		return nil
	}

	reader := csv.NewReader(file)
	const batchSize = 1000
	var batch []village.Village

	// TRY TO READ HEADER
	if _, err := reader.Read(); err != nil {
		return err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading village.csv:", err)
			continue
		}

		if len(record) < 3 {
			continue
		}

		id := record[0]
		districtID := record[2]

		batch = append(batch, village.Village{
			ID:         id,
			Name:       record[1],
			DistrictID: districtID,
		})
	}

	// INSERT DATA IN BATCHES
	if len(batch) > 0 {
		result := db.CreateInBatches(batch, batchSize)

		if result.Error != nil {
			log.Println("✗", result.Error)
			return result.Error
		}

		batch = batch[:0] // RESET BATCH

		log.Println("✓ Village table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
