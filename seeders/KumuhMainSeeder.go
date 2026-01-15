package seeders

import (
	"encoding/csv"
	"klinik-pkp-api/internal/api/kumuh"
	"log"
	"os"
	"gorm.io/gorm"
)

// SEEDS FROM THE CSV FILE
func KumuhMainSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/kumuh.csv")

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
	var data []kumuh.KumuhMain

	db.Model(&kumuh.KumuhMain{}).Count(&count)

	if count > 0 {
		log.Println("!? Kawasan Kumuh table already seeded, skipping...")

		return nil
	}

	// _, records[1:] SKIP HEADER ROW
	for _, record := range records[1:] {
		if len(record) < 2 {
			continue
		}

		// id := record[0]

		kumuh := kumuh.KumuhMain{
			// ID:   id,
			// Name: record[1],
		}

		data = append(data, kumuh)
	}

	if len(data) > 0 {
		result := db.CreateInBatches(data, 100)

		if result.Error != nil {
			return result.Error
		}

		log.Println("✓ Kumuh table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
