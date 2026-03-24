package district

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
type DistrictPayload struct {
	Name     string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
	RegionID string `json:"region_id" validate:"required"`
}

// FILTER FOR QUERY PARAMETERS
type DistrictFilter struct {
	RegionID *uint64
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: utils.NewValidator()}
}

func (s *Service) GetDistricts(filter DistrictFilter) ([]District, error) {
	var districts []District

	query := s.db.Preload("Region.Province").Model(&District{})

	// FILTER BY REGION
	if filter.RegionID != nil {
		query = query.Where("district.region_id = ?", *filter.RegionID)
	}

	if err := query.Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

func (s *Service) GetDistrictById(id uint64) (*District, error) {
	var district District

	if err := s.db.Preload("Region.Province").First(&district, id).Error; err != nil {
		return nil, err
	}

	return &district, nil
}

func (s *Service) AddDistrict(payload DistrictPayload) (*District, error) {
	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return nil, err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	newDistrict := District{
		Name:     payload.Name,
		RegionID: regionID,
	}

	if err := tx.Create(&newDistrict).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newDistrict, nil
}

func (s *Service) EditDistrictById(id uint64, payload DistrictPayload) error {
	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	newDistrict := District{
		Name:     payload.Name,
		RegionID: regionID,
	}

	result := tx.Model(&District{}).Where("id = ?", id).Updates(&newDistrict)

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

// SOFT DELETE
func (s *Service) DeleteDistrictById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&District{}, id)

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
