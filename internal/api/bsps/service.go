package bsps

import (
	"fmt"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type BSPSPayload struct {
	VillageID   string             `json:"village_id"`
	DistrictID  string             `json:"district_id"`
	RegionID    string             `json:"region_id"`
	UnitCount   int                `json:"unit_count"`
	YearGiven   int                `json:"year_given"`
	Coordinates []utils.Coordinate `json:"coordinates"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetBSPS() ([]BSPS, error) {
	var bsps []BSPS

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Find(&bsps).Error; err != nil {
		return nil, err
	}

	return bsps, nil
}

func (s *Service) GetBSPSById(id uint64) (*BSPS, error) {
	var bsps BSPS

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		First(&bsps, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("BSPS with id %d is not found", id)
		} else {
			return nil, err
		}
	}

	return &bsps, nil
}

func (s *Service) AddBSPS(payload BSPSPayload) (*BSPS, error) {
	// VALIDATIONS
	err := s.validator.ValidateStruct(map[string]any{
		"VillageID":   payload.VillageID,
		"DistrictID":  payload.DistrictID,
		"RegionID":    payload.RegionID,
		"UnitCount":   payload.UnitCount,
		"YearGiven":   payload.YearGiven,
		"Coordinates": payload.Coordinates,
	})

	if err != nil {
		return nil, err
	}

	var bsps BSPS

	// ENSURE RELATED VILLAGE EXISTS
	var village village.Village

	if err := s.db.First(&village, "id = ?", payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("village with id %s is not found", payload.VillageID)
		}

		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district district.District

	if err := s.db.First(&district, "id = ?", payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("district with id %s is not found", payload.DistrictID)
		}

		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region region.Region

	if err := s.db.First(&region, "id = ?", payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("region with id %s is not found", payload.RegionID)
		}

		return nil, err
	}

	bsps.VillageID = payload.VillageID
	bsps.DistrictID = payload.DistrictID
	bsps.RegionID = payload.RegionID
	bsps.UnitCount = payload.UnitCount
	bsps.YearGiven = payload.YearGiven
	bsps.Coordinates = payload.Coordinates

	if err := s.db.Create(&bsps).Error; err != nil {
		return nil, err
	}

	return &bsps, nil
}

func (s *Service) EditBSPSById(id uint64, payload BSPSPayload) (*BSPS, error) {
	// VALIDATIONS
	err := s.validator.ValidateStruct(map[string]any{
		"VillageID":   payload.VillageID,
		"DistrictID":  payload.DistrictID,
		"RegionID":    payload.RegionID,
		"UnitCount":   payload.UnitCount,
		"YearGiven":   payload.YearGiven,
		"Coordinates": payload.Coordinates,
	})

	if err != nil {
		return nil, err
	}

	var bsps BSPS

	// ENSURE BSPS EXISTS
	if err := s.db.First(&bsps, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("BSPS with id %d is not found", id)
		}

		return nil, err
	}

	// ENSURE RELATED VILLAGE EXISTS
	var village village.Village

	if err := s.db.First(&village, "id = ?", payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("village with id %s is not found", payload.VillageID)
		}

		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district district.District

	if err := s.db.First(&district, "id = ?", payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("district with id %s is not found", payload.DistrictID)
		}

		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region region.Region

	if err := s.db.First(&region, "id = ?", payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("region with id %s is not found", payload.RegionID)
		}

		return nil, err
	}

	bsps.VillageID = payload.VillageID
	bsps.DistrictID = payload.DistrictID
	bsps.RegionID = payload.RegionID
	bsps.UnitCount = payload.UnitCount
	bsps.YearGiven = payload.YearGiven
	bsps.Coordinates = payload.Coordinates

	if err := s.db.Save(&bsps).Error; err != nil {
		return nil, err
	}

	return &bsps, nil
}

func (s *Service) DeleteBSPSById(id uint64) error {
	var bsps BSPS

	if err := s.db.First(&bsps, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("BSPS with id %d is not found", id)
		}

		return err
	}

	if err := s.db.Delete(&bsps).Error; err != nil {
		return err
	}

	return nil
}
