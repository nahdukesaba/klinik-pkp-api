package seeders

import (
	"encoding/csv"
	"gorm.io/gorm"
	"io"
	"klinik-pkp-api/models"
	"log"
	"os"
	"strconv"
)

// SEEDS FROM PREDEFINED DATA
func RusunSeeder(db *gorm.DB) error {
	file, err := os.Open("script/rusun.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64

	db.Model(&models.Rusun{}).Count(&count)

	if count > 0 {
		log.Println("!? Rusun table already seeded, skipping...")

		return nil
	}

	reader := csv.NewReader(file)
	const batchSize = 100
	var batch []models.Rusun

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

		floor, _ := strconv.Atoi(record[10])
		unitCount, _ := strconv.Atoi(record[11])
		latitude, _ := strconv.ParseFloat(record[6], 64)
		longitude, _ := strconv.ParseFloat(record[7], 64)
		coordinate := models.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		}

		batch = append(batch, models.Rusun{
			Name:       record[0],
			VillageID:  record[1],
			DistrictID: record[2],
			RegionID:   record[3],
			ProvinceID: record[4],
			Address:    record[5],
			Tower:      record[8],
			UnitType:   record[9],
			Floor:      floor,
			UnitCount:  unitCount,
			Coordinate: []models.Coordinate{coordinate},
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
