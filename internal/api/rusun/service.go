package rusun

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/province"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
)

// Service struct
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// RusunPayload represents data from input
type RusunPayload struct {
	VillageID   string             `json:"village_id"`
	DistrictID  string             `json:"district_id"`
	RegionID    string             `json:"region_id"`
	ProvinceID  string             `json:"province_id"`
	Name        string             `json:"name"`
	Address     string             `json:"address"`
	Tower       int                `json:"tower"`
	UnitType    string             `json:"unit_type"`
	Floor       int                `json:"floor"`
	UnitCount   int                `json:"unit_count"`
	YearGiven   int                `json:"year_given"`
	Coordinates []utils.Coordinate `json:"coordinates"`
}

// NewService constructor
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetAll() ([]Rusun, error) {
	var rusuns []Rusun

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Preload("Province").
		Find(&rusuns).Error; err != nil {
		return nil, err
	}

	return rusuns, nil
}

func (s *Service) GetByID(id int) (*Rusun, error) {
	var rusun Rusun
	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Preload("Province").
		First(&rusun, id).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *Service) Create(payload RusunPayload) (*Rusun, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"VillageID":   payload.VillageID,
		"DistrictID":  payload.DistrictID,
		"RegionID":    payload.RegionID,
		"ProvinceID":  payload.ProvinceID,
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

	// ENSURE RELATIONS EXIST
	var vill village.Village

	if err := s.db.First(&vill, payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}
		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var dist district.District

	if err := s.db.First(&dist, payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var reg region.Region

	if err := s.db.First(&reg, payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}
		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var prov province.Province

	if err := s.db.First(&prov, payload.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}
		return nil, err
	}

	rusun := Rusun{
		VillageID:   payload.VillageID,
		DistrictID:  payload.DistrictID,
		RegionID:    payload.RegionID,
		ProvinceID:  payload.ProvinceID,
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

func (s *Service) Update(id int, payload RusunPayload) (*Rusun, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"VillageID":   payload.VillageID,
		"DistrictID":  payload.DistrictID,
		"RegionID":    payload.RegionID,
		"ProvinceID":  payload.ProvinceID,
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
			return nil, errors.New("Rusun not found")
		}

		return nil, err
	}

	// ENSURE RELATED VILLAGE EXISTS
	var vill village.Village

	if err := s.db.First(&vill, payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}
		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var dist district.District

	if err := s.db.First(&dist, payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var reg region.Region

	if err := s.db.First(&reg, payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var prov province.Province

	if err := s.db.First(&prov, payload.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	rusun.VillageID = payload.VillageID
	rusun.DistrictID = payload.DistrictID
	rusun.RegionID = payload.RegionID
	rusun.ProvinceID = payload.ProvinceID
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

func (s *Service) Delete(id int) error {
	var rusun Rusun

	if err := s.db.First(&rusun, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Rusun not found")
		}

		return err
	}

	if err := s.db.Delete(&rusun).Error; err != nil {
		return err
	}

	return nil
}
