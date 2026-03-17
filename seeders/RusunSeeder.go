package seeders

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/internal/api/v1/rusun"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func RusunSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/rusun.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&rusun.Rusun{}).Count(&count)

	if count > 10 {
		log.Println("!? Rusun table already seeded, skipping...")

		return nil
	}

	var batch []rusun.Rusun
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
			return fmt.Errorf("Row %d: Error reading rusun.csv: %v", rowNum, err)
		}

		if len(record) < 12 {
			skippedRows++
			continue
		}

		isVillageIdNumeric, _ := regexp.MatchString(`^\d+$`, strings.TrimSpace(record[0]))
		if isVillageIdNumeric == false {
			log.Printf("Row %d: Skipping - Village ID is not a valid unsigned integer (%q)\n", rowNum, record[0])
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

		tower, err := strconv.ParseUint(strings.TrimSpace(record[5]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Tower (%q): %v", rowNum, record[5], err)
		}

		floor, err := strconv.ParseUint(strings.TrimSpace(record[7]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Floor (%q): %v", rowNum, record[7], err)
		}

		unitCount, err := strconv.ParseUint(strings.TrimSpace(record[8]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Unit Count (%q): %v", rowNum, record[8], err)
		}

		yearGiven, err := strconv.ParseUint(strings.TrimSpace(record[9]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Year Given (%q): %v", rowNum, record[9], err)
		}

		var coordinate utils.Coordinate

		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[10]), 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing coordinates: %v", rowNum, err)
		}

		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[11]), 64)
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

		batch = append(batch, rusun.Rusun{
			VillageID:  villageID,
			DistrictID: districtID,
			RegionID:   regionID,
			Name:       record[3],
			Address:    record[4],
			Tower:      tower,
			UnitType:   record[6],
			Floor:      floor,
			UnitCount:  unitCount,
			YearGiven:  yearGiven,
			Coordinate: coordinate,
		})
	}

	log.Printf("Rusun: %d non-data rows skipped\n", skippedRows)

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		batch = batch[:0] // RESET BATCH

		log.Println("✓ Rusun table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
