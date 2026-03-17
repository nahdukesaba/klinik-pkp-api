package seeders

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/internal/api/v1/kumuh"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func KumuhSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/kumuh.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&kumuh.KawasanKumuh{}).Count(&count)

	if count > 10 {
		log.Println("!? Kawasan Kumuh table already seeded, skipping...")

		return nil
	}

	var batch []kumuh.KawasanKumuh
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
			return fmt.Errorf("Row %d: Error reading kumuh.csv: %v", rowNum, err)
		}

		if len(record) < 11 {
			skippedRows++
			continue
		}

		// SKIP ROW IF DISTRICT ID IS NOT NUMERIC
		isDistrictIdNumeric, _ := regexp.MatchString(`^\d+$`, strings.TrimSpace(record[3]))
		if isDistrictIdNumeric == false {
			log.Printf("Row %d: Skipping - District ID is not numeric (%q)\n", rowNum, record[3])
			skippedRows++
			continue
		}

		districtId, err := strconv.ParseUint(strings.TrimSpace(record[3]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing District ID (%q): %v", rowNum, record[3], err)
		}

		regionId, err := strconv.ParseUint(strings.TrimSpace(record[4]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Region ID (%q): %v", rowNum, record[4], err)
		}

		totalArea, err := strconv.ParseFloat(strings.TrimSpace(record[5]), 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Total Area (%q): %v", rowNum, record[5], err)
		}

		totalPopulation, err := strconv.ParseUint(strings.TrimSpace(record[6]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Total Population (%q): %v", rowNum, record[6], err)
		}

		slumValue, err := strconv.ParseUint(strings.TrimSpace(record[7]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Slum Value (%q): %v", rowNum, record[7], err)
		}

		yearInspected, err := strconv.ParseUint(strings.TrimSpace(record[8]), 10, 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing Year Inspected (%q): %v", rowNum, record[8], err)
		}

		var coordinate utils.Coordinate

		latitude, err := strconv.ParseFloat(strings.TrimSpace(record[9]), 64)
		if err != nil {
			return fmt.Errorf("Row %d: Error parsing coordinates: %v", rowNum, err)
		}

		longitude, err := strconv.ParseFloat(strings.TrimSpace(record[10]), 64)
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

		batch = append(batch, kumuh.KawasanKumuh{
			DistrictID:      districtId,
			RegionID:        regionId,
			AreaName:        record[0],
			Environments:    record[1],
			Villages:        record[2],
			TotalArea:       totalArea,
			TotalPopulation: totalPopulation,
			SlumValue:       slumValue,
			YearInspected:   yearInspected,
			Coordinate:      coordinate,
		})
	}

	log.Printf("Kawasan Kumuh: %d non-data rows skipped\n", skippedRows)

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		batch = batch[:0] // RESET BATCH

		log.Println("✓ Kawasan Kumuh table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
