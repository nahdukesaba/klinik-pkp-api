package seeders

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/internal/api/bsps"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

// SEEDS FROM PREDEFINED DATA
func BSPSSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/bsps.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&bsps.BSPS{}).Count(&count)

	if count > 0 {
		log.Println("!? BSPS table already seeded, skipping...")

		return nil
	}

	reader := csv.NewReader(file)
	const batchSize = 100
	var batch []bsps.BSPS

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
			log.Println("Error reading bsps.csv:", err)
			continue
		}

		if len(record) < 8 {
			continue
		}

		villageID, err := strconv.ParseUint(record[0], 10, 64)
		districtID, err := strconv.ParseUint(record[1], 10, 64)
		regionID, err := strconv.ParseUint(record[2], 10, 64)
		unitCount, err := strconv.ParseUint(record[3], 10, 64)
		yearGiven, err := strconv.ParseUint(record[4], 10, 64)

		var coordinate utils.Coordinate

		// PARSE COORDINATES FROM COLUMN[6] ONWARDS
		for i := 6; i+1 < len(record); i += 2 {
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

			coordinate = utils.Coordinate{
				Latitude:  latitude,
				Longitude: longitude,
			}
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

	if len(batch) > 0 {
		result := db.CreateInBatches(batch, 100)

		if result.Error != nil {
			return result.Error
		}

		// RESET SEQUENCE AFTER MANUAL ID INSERTS (OPTIONAL, BEST TO LET DATABASE HANDLE IDS)
		// if err := db.Exec(`
		// 	SELECT setval(
		// 		pg_get_serial_sequence('penerimaan_bsps', 'id'),
		// 		(SELECT COALESCE(MAX(id), 1) FROM penerimaan_bsps),
		// 		true
		// 	)
		// `).Error; err != nil {
		// 	return err
		// }

		batch = batch[:0] // RESET BATCH

		log.Println("✓ BSPS table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
