package bsps

import (
	"klinik-pkp-api/utils"
	"gorm.io/gorm"
)

// Service struct
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// BSPSPayload represents data from input
type BSPSPayload struct {
	VillageID            string   `json:"village_id"`
	DistrictID           string   `json:"district_id"`
	RegionID             string   `json:"region_id"`
	ProvinceID           string   `json:"province_id"`
	Name                 string   `json:"name"`
	Address              string   `json:"address"`
	MGRSCoordinate       string   `json:"mgrs_coordinate"`
	GeographicCoordinate string   `json:"geographic_coordinate"`
	DecimalCoordinate    []string `json:"decimal_coordinate"`
	Tower                string   `json:"tower"`
	UnitType             string   `json:"unit_type"`
	Floor                int      `json:"floor"`
	UnitCount            int      `json:"unit_count"`
}

// NewService constructor
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetAll() ([]BSPS, error) {
	var bsps []BSPS

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Preload("Province").
		Find(&bsps).Error; err != nil {
		return nil, err
	}

	return bsps, nil
}
