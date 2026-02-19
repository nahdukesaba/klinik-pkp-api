package seeders

import (
	"encoding/csv"
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/kumuh"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

// SEEDS FROM THE CSV FILE
func KawasanKumuhSeeder(db *gorm.DB) error {
	file, err := os.Open("seeders/data/kawasan-kumuh.csv")

	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return err
	}

	// SKIP SEED IF TABLE IS NOT EMPTY
	var count int64
	var data []kumuh.KawasanKumuh

	db.Model(&kumuh.KawasanKumuh{}).Count(&count)

	if count > 10 {
		log.Println("!? Kawasan Kumuh table already seeded, skipping...")

		return nil
	}

	// _, records[1:] SKIP HEADER ROW
	for _, record := range records[1:] {
		if len(record) < 10 {
			continue
		}

		districtId, err := strconv.ParseUint(record[3], 10, 64)

		if err != nil {
			log.Println("Error parsing District ID:", err)

			continue
		}

		regionId, err := strconv.ParseUint(record[4], 10, 64)

		if err != nil {
			log.Println("Error parsing Region ID:", err)

			continue
		}

		totalArea, err := strconv.ParseFloat(record[5], 64)

		if err != nil {
			log.Println("Error parsing Total Area:", err)

			continue
		}

		totalPopulation, err := strconv.ParseUint(record[6], 10, 64)

		if err != nil {
			log.Println("Error parsing Total Population:", err)

			continue
		}

		slumValue, err := strconv.ParseUint(record[7], 10, 64)

		if err != nil {
			log.Println("Error parsing Slum Value:", err)

			continue
		}

		var coordinate utils.Coordinate

		stringLatitude := strings.TrimSpace(record[8])
		stringLongitude := strings.TrimSpace(record[9])

		if stringLatitude == "" || stringLongitude == "" {
			continue
		}

		latitude, latErr := strconv.ParseFloat(stringLatitude, 64)
		longitude, longErr := strconv.ParseFloat(stringLongitude, 64)

		if latErr != nil || longErr != nil {
			return fmt.Errorf("invalid coordinates format: %v %v", latErr, longErr)
		}

		utils.ValidateCoordinate(utils.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		})

		coordinate = utils.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		}

		kumuh := kumuh.KawasanKumuh{
			AreaName:        record[0],
			Environments:    record[1],
			Villages:        record[2],
			DistrictID:      districtId,
			RegionID:        regionId,
			TotalArea:       totalArea,
			TotalPopulation: totalPopulation,
			SlumValue:       slumValue,
			Coordinate:      coordinate,
		}

		data = append(data, kumuh)
	}

	if len(data) > 0 {
		result := db.CreateInBatches(data, 100)

		if result.Error != nil {
			return result.Error
		}

		log.Println("✓ Kawasan Kumuh table seeded successfully with", result.RowsAffected, "records")
	}

	return nil
}
