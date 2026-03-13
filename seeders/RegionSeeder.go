package seeders

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/internal/api/region"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// SEEDS FROM THE CSV FILE
func RegionSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/region.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&region.Region{}).Count(&count)

	if count > 10 {
		log.Println("!? Kawasan Kumuh table already seeded, skipping...")

		return nil
	}

	var batch []region.Region
	reader := csv.NewReader(file)
	rowNum := 0
	skippedRows := 0

	for {
		record, err := reader.Read()
		rowNum++

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("Row %d: Error reading region.csv: %v", rowNum, err)
		}

		if len(record) < 3 {
			skippedRows++
			continue
		}

		isIdNumeric, _ := regexp.MatchString(`^\d+$`, strings.TrimSpace(record[0]))
		if isIdNumeric == false {
			log.Printf("Row %d: Skipping - ID is not a valid unsigned integer (%q)\n", rowNum, record[0])
			skippedRows++
			continue
		}

		id, err := strconv.ParseUint(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing ID (%q): %v", rowNum, record[0], err)
		}

		provinceID, err := strconv.ParseUint(strings.TrimSpace(record[2]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing ProvinceID (%q): %v", rowNum, record[2], err)
		}

		batch = append(batch, region.Region{
			ID:         id,
			Name:       record[1],
			ProvinceID: provinceID,
		})
	}

	log.Printf("Region: %d non-data rows skipped\n", skippedRows)

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		log.Println("✓ Region table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
