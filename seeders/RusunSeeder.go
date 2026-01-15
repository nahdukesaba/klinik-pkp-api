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
	file, err := os.Open("seeders/data/rusun.csv")

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

		if len(record) < 12 {
			continue
		}

		tower, _ := strconv.Atoi(record[5])
		floor, _ := strconv.Atoi(record[7])
		unitCount, _ := strconv.Atoi(record[8])
		yearGiven, _ := strconv.Atoi(record[9])

		var coordinates []utils.Coordinate

		// PARSE COORDINATES FROM COLUMN[10] ONWARDS
		for i := 10; i+1 < len(record); i += 2 {
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

			utils.ValidateCoordinate(utils.Coordinate{
				Latitude:  latitude,
				Longitude: longitude,
			})

			coordinates = append(coordinates, utils.Coordinate{
				Latitude:  latitude,
				Longitude: longitude,
			})
		}

		batch = append(batch, rusun.Rusun{
			VillageID:   record[0],
			DistrictID:  record[1],
			RegionID:    record[2],
			Name:        record[3],
			Address:     record[4],
			Tower:       tower,
			UnitType:    record[6],
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
