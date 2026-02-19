package bsps

import (
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
	"strconv"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type BSPSPayload struct {
	VillageID   string             `json:"village_id" validate:"required"`
	DistrictID  string             `json:"district_id" validate:"required"`
	RegionID    string             `json:"region_id" validate:"required"`
	UnitCount   uint64             `json:"unit_count"`
	YearGiven   uint64             `json:"year_given"`
	Status      string             `json:"status" validate:"required,ne=Null,ne=null,ne=NULL"`
	Coordinate  utils.Coordinate   `json:"coordinate" validate:"required" gorm:"-"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetBSPS() ([]BSPS, error) {
	var bsps []BSPS

	if err := s.db.Preload("Village").Preload("District").Preload("Region").Find(&bsps).Error; err != nil {
		return nil, err
	}

	return bsps, nil
}

func (s *Service) GetBSPSById(id uint64) (*BSPS, error) {
	var bsps BSPS

	if err := s.db.Preload("Village").Preload("District").Preload("Region.Province").First(&bsps, id).Error; err != nil {
		return nil, err
	}

	return &bsps, nil
}

func (s *Service) AddBSPS(payload *BSPSPayload) (*BSPS, error) {
	villageID, err := strconv.ParseUint(payload.VillageID, 10, 64)

	if err != nil {
		return nil, err
	}

	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return nil, err
	}

	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return nil, err
	}

	newBSPS := BSPS{
		VillageID:   villageID,
		DistrictID:  districtID,
		RegionID:    regionID,
		UnitCount:   payload.UnitCount,
		YearGiven:   payload.YearGiven,
		Status:      payload.Status,
		Coordinate: payload.Coordinate,
	}

	if err := s.db.Create(&newBSPS).Error; err != nil {
		return nil, err
	}

	return &newBSPS, nil
}

func (s *Service) EditBSPSById(id uint64, payload *BSPSPayload) error {
	villageID, err := strconv.ParseUint(payload.VillageID, 10, 64)

	if err != nil {
		return err
	}

	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return err
	}

	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return err
	}

	// CHECK IF RECORD EXISTS
	existingBSPS := BSPS{}

	if err := s.db.First(&existingBSPS, id).Error; err != nil {
		return err
	}

	newBSPS := BSPS{
		VillageID:   villageID,
		DistrictID:  districtID,
		RegionID:    regionID,
		UnitCount:   payload.UnitCount,
		YearGiven:   payload.YearGiven,
		Status:      payload.Status,
		Coordinate: payload.Coordinate,
	}

	result := s.db.Model(&BSPS{}).Where("id = ?", id).Updates(newBSPS)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// SOFT DELETE
func (s *Service) DeleteBSPSById(id uint64) error {
	result := s.db.Delete(&BSPS{}, id)

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
