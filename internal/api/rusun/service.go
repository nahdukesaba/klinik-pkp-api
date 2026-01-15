package rusun

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type RusunPayload struct {
	VillageID   string             `json:"village_id"`
	DistrictID  string             `json:"district_id"`
	RegionID    string             `json:"region_id"`
	Name        string             `json:"name"`
	Address     string             `json:"address"`
	Tower       int                `json:"tower"`
	UnitType    string             `json:"unit_type"`
	Floor       int                `json:"floor"`
	UnitCount   int                `json:"unit_count"`
	YearGiven   int                `json:"year_given"`
	Coordinates []utils.Coordinate `json:"coordinates"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetRusun() ([]Rusun, error) {
	var rusuns []Rusun

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Find(&rusuns).Error; err != nil {
		return nil, err
	}

	return rusuns, nil
}

func (s *Service) GetRusunById(id uint64) (*Rusun, error) {
	var rusun Rusun

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		First(&rusun, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("rusun with id %d is not found", id)
		}
		return nil, err
	}

	return &rusun, nil
}

func (s *Service) AddRusun(payload RusunPayload) (*Rusun, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"VillageID":   payload.VillageID,
		"DistrictID":  payload.DistrictID,
		"RegionID":    payload.RegionID,
		"Name":        payload.Name,
		"Address":     payload.Address,
		"Tower":       payload.Tower,
		"UnitType":    payload.UnitType,
		"Floor":       payload.Floor,
		"UnitCount":   payload.UnitCount,
		"YearGiven":   payload.YearGiven,
		"Coordinates": payload.Coordinates,
	})

	if err != nil {
		return nil, err
	}

	// ENSURE RELATED VILLAGE EXISTS
	var village village.Village

	if err := s.db.First(&village, payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("village with id %s is not found", payload.VillageID)
		}

		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district district.District

	if err := s.db.First(&district, payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("district with id %s is not found", payload.DistrictID)
		}
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region region.Region

	if err := s.db.First(&region, payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("region with id %s is not found", payload.RegionID)
		}

		return nil, err
	}

	rusun := Rusun{
		VillageID:   payload.VillageID,
		DistrictID:  payload.DistrictID,
		RegionID:    payload.RegionID,
		Name:        payload.Name,
		Address:     payload.Address,
		Tower:       payload.Tower,
		UnitType:    payload.UnitType,
		Floor:       payload.Floor,
		UnitCount:   payload.UnitCount,
		YearGiven:   payload.YearGiven,
		Coordinates: payload.Coordinates,
	}

	if err := s.db.Create(&rusun).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *Service) EditRusunById(id uint64, payload RusunPayload) (*Rusun, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"VillageID":   payload.VillageID,
		"DistrictID":  payload.DistrictID,
		"RegionID":    payload.RegionID,
		"Name":        payload.Name,
		"Address":     payload.Address,
		"Tower":       payload.Tower,
		"UnitType":    payload.UnitType,
		"Floor":       payload.Floor,
		"UnitCount":   payload.UnitCount,
		"YearGiven":   payload.YearGiven,
		"Coordinates": payload.Coordinates,
	})

	if err != nil {
		return nil, err
	}

	var rusun Rusun

	// ENSURE RUSUN EXISTS
	if err := s.db.First(&rusun, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("rusun with id %d is not found", id)
		}

		return nil, err
	}

	// ENSURE RELATED VILLAGE EXISTS
	var village village.Village

	if err := s.db.First(&village, payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("village with id %s is not found", payload.VillageID)
		}
		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district district.District

	if err := s.db.First(&district, payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("district with id %s is not found", payload.DistrictID)
		}
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region region.Region

	if err := s.db.First(&region, payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("region with id %s is not found", payload.RegionID)
		}

		return nil, err
	}

	rusun.VillageID = payload.VillageID
	rusun.DistrictID = payload.DistrictID
	rusun.RegionID = payload.RegionID
	rusun.Name = payload.Name
	rusun.Address = payload.Address
	rusun.Coordinates = payload.Coordinates
	rusun.Tower = payload.Tower
	rusun.UnitType = payload.UnitType
	rusun.Floor = payload.Floor
	rusun.UnitCount = payload.UnitCount

	if err := s.db.Save(&rusun).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *Service) DeleteRusunById(id uint64) error {
	var rusun Rusun

	if err := s.db.First(&rusun, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("rusun with id %d is not found", id)
		}

		return err
	}

	if err := s.db.Delete(&rusun).Error; err != nil {
		return err
	}

	return nil
}
