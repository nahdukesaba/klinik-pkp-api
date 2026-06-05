package seeders

import (
	"encoding/csv"
	"fmt"
	"io"
	"klinik-pkp-api/internal/api/v1/bsps"
	"klinik-pkp-api/internal/api/v1/sosialisasi"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func SosialisasiSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/sosialisasi.csv")

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

	var batch []sosialisasi.Sosialisasi
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
			return fmt.Errorf("Row %d: Error reading sosialisasi.csv: %v", rowNum, err)
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

		batch = append(batch, sosialisasi.Sosialisasi{
			VillageID:   villageID,
			DistrictID:  districtID,
			RegionID:    regionID,
			Title:       record[3],
			Location:    record[4],
			Description: record[5],
			Coordinate:  coordinate,
		})
	}

	log.Printf("Sosialisasi: %d non-data rows skipped\n", skippedRows)

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		batch = batch[:0] // RESET BATCH

		log.Println("✓ Sosialisasi table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
