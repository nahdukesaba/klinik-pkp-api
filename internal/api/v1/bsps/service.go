package bsps

import (
	"klinik-pkp-api/utils"
	"strconv"

	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type BSPSPayload struct {
	VillageID  string           `json:"village_id" validate:"required"`
	DistrictID string           `json:"district_id" validate:"required"`
	RegionID   string           `json:"region_id" validate:"required"`
	UnitCount  uint64           `json:"unit_count"`
	YearGiven  uint64           `json:"year_given"`
	Status     string           `json:"status" validate:"required,ne=Null,ne=null,ne=NULL"`
	Coordinate utils.Coordinate `json:"coordinate" validate:"required" gorm:"-"`
}

// FILTER FOR QUERY PARAMETERS
type BSPSFilter struct {
	VillageID  *uint64
	DistrictID *uint64
	RegionID   *uint64
	Status     *string
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: utils.NewValidator()}
}

func (s *Service) GetBSPS(filter BSPSFilter) ([]BSPS, error) {
	var bsps []BSPS

	query := s.db.Preload("Village").Preload("District").Preload("Region").Model(&BSPS{})

	// FILTER BY VILLAGE
	if filter.VillageID != nil {
		query = query.Where("bsps.village_id = ?", *filter.VillageID)
	}

	// FILTER BY DISTRICT
	if filter.DistrictID != nil {
		query = query.Where("bsps.district_id = ?", *filter.DistrictID)
	}

	// FILTER BY REGION
	if filter.RegionID != nil {
		query = query.Where("bsps.region_id = ?", *filter.RegionID)
	}

	// FILTER BY STATUS
	if filter.Status != nil {
		query = query.Where("bsps.status = ?", *filter.Status)
	}

	if err := query.Find(&bsps).Error; err != nil {
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

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	newBSPS := BSPS{
		VillageID:  villageID,
		DistrictID: districtID,
		RegionID:   regionID,
		UnitCount:  payload.UnitCount,
		YearGiven:  payload.YearGiven,
		Status:     payload.Status,
		Coordinate: payload.Coordinate,
	}

	if err := tx.Create(&newBSPS).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
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

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// CHECK IF RECORD EXISTS
	existingBSPS := BSPS{}

	if err := tx.First(&existingBSPS, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	newBSPS := BSPS{
		VillageID:  villageID,
		DistrictID: districtID,
		RegionID:   regionID,
		UnitCount:  payload.UnitCount,
		YearGiven:  payload.YearGiven,
		Status:     payload.Status,
		Coordinate: payload.Coordinate,
	}

	result := tx.Model(&BSPS{}).Where("id = ?", id).Updates(newBSPS)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// SOFT DELETE
func (s *Service) DeleteBSPSById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&BSPS{}, id)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
