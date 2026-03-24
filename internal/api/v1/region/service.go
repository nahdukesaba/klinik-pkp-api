package region

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
type RegionPayload struct {
	Name       string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
	ProvinceID string `json:"province_id" validate:"required"`
}

// FILTER FOR QUERY PARAMETERS
type RegionFilter struct {
	ProvinceID *uint64
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: utils.NewValidator()}
}

func (s *Service) GetRegions(filter RegionFilter) ([]Region, error) {
	var regions []Region

	query := s.db.Preload("Province")

	// FILTER BY PROVINCE
	if filter.ProvinceID != nil {
		query = query.Where("province_id = ?", *filter.ProvinceID)
	}

	if err := query.Find(&regions).Error; err != nil {
		return nil, err
	}

	return regions, nil
}

func (s *Service) GetRegionById(id uint64) (*Region, error) {
	var region Region

	if err := s.db.Preload("Province").First(&region, id).Error; err != nil {
		return nil, err
	}

	return &region, nil
}

func (s *Service) AddRegion(payload RegionPayload) (*Region, error) {
	provinceID, err := strconv.ParseUint(payload.ProvinceID, 10, 64)

	if err != nil {
		return nil, err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	newRegion := Region{
		Name:       payload.Name,
		ProvinceID: provinceID,
	}

	if err := tx.Create(&newRegion).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newRegion, nil
}

func (s *Service) EditRegionById(id uint64, payload RegionPayload) error {
	provinceID, err := strconv.ParseUint(payload.ProvinceID, 10, 64)

	if err != nil {
		return err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	newRegion := Region{
		Name:       payload.Name,
		ProvinceID: provinceID,
	}

	result := tx.Model(&Region{}).Where("id = ?", id).Updates(&newRegion)

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
func (s *Service) DeleteRegionById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&Region{}, id)

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
