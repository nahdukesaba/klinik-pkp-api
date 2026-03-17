package seeders

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/internal/api/v1/bsps"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func BSPSSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/bsps.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&bsps.BSPS{}).Count(&count)

	if count > 10 {
		log.Println("!? BSPS table already seeded, skipping...")

		return nil
	}

	var batch []bsps.BSPS
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
			return fmt.Errorf("Row %d: Error reading bsps.csv: %v", rowNum, err)
		}

		if len(record) < 8 {
			skippedRows++
			continue
		}

		isVillageIdNumeric, _ := regexp.MatchString(`^\d+$`, strings.TrimSpace(record[0]))
		if isVillageIdNumeric == false {
			log.Printf("Row %d: Skipping - Village ID is not numeric (%q)\n", rowNum, record[0])
			skippedRows++
			continue
		}

		villageID, err := strconv.ParseUint(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Village ID (%q): %v", rowNum, record[0], err)
		}

		districtID, err := strconv.ParseUint(strings.TrimSpace(record[1]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing District ID (%q): %v", rowNum, record[1], err)
		}

		regionID, err := strconv.ParseUint(strings.TrimSpace(record[2]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Region ID (%q): %v", rowNum, record[2], err)
		}

		unitCount, err := strconv.ParseUint(strings.TrimSpace(record[3]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Unit Count (%q): %v", rowNum, record[3], err)
		}

		yearGiven, err := strconv.ParseUint(strings.TrimSpace(record[4]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Year Given (%q): %v", rowNum, record[4], err)
		}

		var coordinate utils.Coordinate

		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[6]), 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing coordinates: %v", rowNum, err)
		}

		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[7]), 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing coordinates: %v", rowNum, err)
		}

		utils.ValidateCoordinate(utils.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		})

		coordinate = utils.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		}

		batch = append(batch, bsps.BSPS{
			VillageID:  villageID,
			DistrictID: districtID,
			RegionID:   regionID,
			UnitCount:  unitCount,
			YearGiven:  yearGiven,
			Status:     record[5],
			Coordinate: coordinate,
		})
	}

	log.Printf("BSPS: %d non-data rows skipped\n", skippedRows)

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		batch = batch[:0] // RESET BATCH

		log.Println("✓ BSPS table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
