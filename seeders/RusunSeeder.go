package seeders

import (
	"encoding/csv"
	"fmt"
	"io"
	"klinik-pkp-api/internal/api/rusun"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// SEEDS FROM PREDEFINED DATA
func RusunSeeder(db *gorm.DB) error {
	file, err := os.Open("others/rusun.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&rusun.Rusun{}).Count(&count)

	if count > 0 {
		log.Println("!? Rusun table already seeded, skipping...")

		return nil
	}

	reader := csv.NewReader(file)
	const batchSize = 100
	var batch []rusun.Rusun

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
			log.Println("Error reading rusun.csv:", err)
			continue
		}

		if len(record) < 9 {
			continue
		}

		id, _ := strconv.ParseUint(record[0], 10, 64)
		tower, _ := strconv.Atoi(record[7])
		floor, _ := strconv.Atoi(record[9])
		unitCount, _ := strconv.Atoi(record[10])
		yearGiven, _ := strconv.Atoi(record[11])

		var coordinates []utils.Coordinate

		// PARSE COORDINATES FROM COLUMN[12] ONWARDS
		for i := 12; i+1 < len(record); i += 2 {
			stringLatitude := strings.TrimSpace(record[i])
			stringLongitude := strings.TrimSpace(record[i+1])

			if stringLatitude == "" || stringLongitude == "" {
				continue
			}

			latitude, err1 := strconv.ParseFloat(stringLatitude, 64)
			longitude, err2 := strconv.ParseFloat(stringLongitude, 64)

			if err1 != nil || err2 != nil {
				return fmt.Errorf("invalid coordinate at column %d", i)
			}

			coordinates = append(coordinates, utils.Coordinate{
				Latitude:  latitude,
				Longitude: longitude,
			})
		}

		batch = append(batch, rusun.Rusun{
			ID:          uint(id),
			VillageID:   record[1],
			DistrictID:  record[2],
			RegionID:    record[3],
			ProvinceID:  record[4],
			Name:        record[5],
			Address:     record[6],
			Tower:       tower,
			UnitType:    record[8],
			Floor:       floor,
			UnitCount:   unitCount,
			YearGiven:   yearGiven,
			Coordinates: coordinates,
		})
	}

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
