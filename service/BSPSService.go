package service

import (
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"
)

// BSPS SERVICE STRUCT
type BSPSService struct {
	db        *gorm.DB
	validator *utils.Validator
}

// REPRESENTS DATA FROM INPUT
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

// RUSUN SERVICE CONSTRUCTOR
func NewBSPSService(db *gorm.DB) *BSPSService {
	return &BSPSService{db: db, validator: &utils.Validator{}}
}

func (s *BSPSService) GetAll() ([]models.BSPS, error) {
	var bsps []models.BSPS

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
