package kumuh

import (
	"fmt"
	"klinik-pkp-api/internal/api/uploads"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
	"strconv"

	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db            *gorm.DB
	validator     *utils.Validator
	uploadService *uploads.Service
}

// PAYLOAD FROM REQUEST BODY
type KawasanKumuhPayload struct {
	DistrictID      string           `json:"district_id" validate:"required"`
	RegionID        string           `json:"region_id" validate:"required"`
	AreaName        string           `json:"area_name" validate:"required,ne=Null,ne=null,ne=NULL"`
	Environments    string           `json:"environments" validate:"required,ne=Null,ne=null,ne=NULL"`
	Villages        string           `json:"villages" validate:"required"`
	TotalArea       float64          `json:"total_area" validate:"required"`
	TotalPopulation uint64           `json:"total_population" validate:"required"`
	SlumValue       uint64           `json:"slum_value" validate:"required"`
	Coordinate      utils.Coordinate `json:"coordinate" validate:"required" gorm:"-"`
}

// FILTER FOR QUERY PARAMETERS
type KawasanKumuhFilter struct {
	RegionID      *uint64
	DistrictID    *uint64
	VillageID     *uint64
	AreaName      *string
	YearInspected *uint64
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		validator: utils.NewValidator(),
	}
}

func (s *Service) GetKumuh(filter KawasanKumuhFilter) ([]KawasanKumuh, error) {
	var kumuh []KawasanKumuh

	query := s.db.Preload("District").Preload("Region").Model(&KawasanKumuh{})

	// FILTER BY REGION
	if filter.RegionID != nil {
		query = query.Where("kawasan_kumuh.region_id = ?", *filter.RegionID)
	}

	// FILTER BY DISTRICT
	if filter.DistrictID != nil {
		query = query.Where("kawasan_kumuh.district_id = ?", *filter.DistrictID)
	}

	// FILTER BY VILLAGE (special case: lookup village name then LIKE match)
	if filter.VillageID != nil {
		var village village.Village

		if err := s.db.First(&village, *filter.VillageID).Error; err != nil {
			return nil, fmt.Errorf("village with id %d not found", *filter.VillageID)
		}

		query = query.Where("kawasan_kumuh.villages LIKE ?", "%"+village.Name+"%")
	}

	// FILTER BY AREA NAME
	if filter.AreaName != nil {
		query = query.Where("kawasan_kumuh.area_name LIKE ?", "%"+*filter.AreaName+"%")
	}

	// FILTER BY YEAR INSPECTED
	if filter.YearInspected != nil {
		query = query.Where("kawasan_kumuh.year_inspected = ?", *filter.YearInspected)
	}

	if err := query.Find(&kumuh).Error; err != nil {
		return nil, err
	}

	return kumuh, nil
}

func (s *Service) GetKumuhById(id uint64) (*KawasanKumuh, error) {
	var kumuh KawasanKumuh

	if err := s.db.Preload("District").Preload("Region").First(&kumuh, id).Error; err != nil {
		return nil, err
	}

	return &kumuh, nil
}

func (s *Service) AddKumuh(payload *KawasanKumuhPayload) (*KawasanKumuh, error) {
	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return nil, err
	}

	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return nil, err
	}

	// CREATE RECORD FIRST TO GET ID
	newKumuh := KawasanKumuh{
		Environments:    payload.Environments,
		AreaName:        payload.AreaName,
		Villages:        payload.Villages,
		DistrictID:      districtID,
		RegionID:        regionID,
		TotalArea:       payload.TotalArea,
		TotalPopulation: payload.TotalPopulation,
		SlumValue:       payload.SlumValue,
		Coordinate:      payload.Coordinate,
	}

	if err := s.db.Create(&newKumuh).Error; err != nil {
		return nil, err
	}

	return &newKumuh, nil
}

func (s *Service) EditKumuhById(id uint64, payload *KawasanKumuhPayload) error {
	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return err
	}

	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return err
	}

	// CHECK IF RECORD EXISTS
	existingKumuh := KawasanKumuh{}

	if err := s.db.First(&existingKumuh, id).Error; err != nil {
		return err
	}

	newKumuh := KawasanKumuh{
		Environments:    payload.Environments,
		AreaName:        payload.AreaName,
		Villages:        payload.Villages,
		DistrictID:      districtID,
		RegionID:        regionID,
		TotalArea:       payload.TotalArea,
		TotalPopulation: payload.TotalPopulation,
		SlumValue:       payload.SlumValue,
		Coordinate:      payload.Coordinate,
	}

	result := s.db.Model(&KawasanKumuh{}).Where("id = ?", id).Updates(newKumuh)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Service) DeleteKumuhById(id uint64) error {
	result := s.db.Delete(&KawasanKumuh{}, id)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
