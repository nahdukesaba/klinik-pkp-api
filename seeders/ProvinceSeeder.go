package seeders

import (
	"encoding/csv"
	"fmt"
	"io"
	"klinik-pkp-api/internal/api/v1/province"
	"log"
	"os"
	"strconv"
	"strings"
	"regexp"
	"gorm.io/gorm"
)

// SEEDS FROM THE CSV FILE
func ProvinceSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/province.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&province.Province{}).Count(&count)

	if count > 10 {
		log.Println("!? Province table already seeded, skipping...")

		return nil
	}

	var batch []province.Province
	reader := csv.NewReader(file)
	rowNum := 0
	skippedRows := 0

	// LOOPS UNTIL EOF
	for {
		record, err := reader.Read()
		rowNum++

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("Row %d: Error reading province.csv: %v", rowNum, err)
		}

		if len(record) < 2 {
			skippedRows++
			continue
		}

		isIdNumeric, _ := regexp.MatchString(`^\d+$`, strings.TrimSpace(record[0]))
		if isIdNumeric == false {
			log.Printf("Row %d: Skipping - ID is not a valid unsigned integer(%q)\n", rowNum, record[0])
			skippedRows++
			continue
		}

		id, err := strconv.ParseUint(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing ID (%q): %v", rowNum, record[0], err)
		}

		batch = append(batch, province.Province{
			ID:   id,
			Name: strings.TrimSpace(record[1]),
		})
	}

	log.Printf("Province: %d non-data rows skipped\n", skippedRows)

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		log.Println("✓ Province table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
