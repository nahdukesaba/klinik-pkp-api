package region

import (
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// POST REGION PAYLOAD FROM REQUEST BODY
type PostRegionPayload struct {
	ID         string `json:"id" validate:"required,min=4,max=4"`
	Name       string `json:"name" validate:"required"`
	ProvinceID string `json:"province_id" validate:"required"`
}

// PUT REGION PAYLOAD FROM REQUEST BODY
type PutRegionPayload struct {
	Name       string `json:"name" validate:"required"`
	ProvinceID string `json:"province_id" validate:"required"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetRegions() ([]Region, error) {
	var regions []Region

	if err := s.db.Preload("Province").Find(&regions).Error; err != nil {
		return nil, err
	}

	return regions, nil
}

func (s *Service) GetRegionById(id string) (*Region, error) {
	var region Region

	if err := s.db.Preload("Province").First(&region, id).Error; err != nil {
		return nil, err
	}

	return &region, nil
}

func (s *Service) AddRegion(payload PostRegionPayload) (*Region, error) {
	region := Region{
		ID:         payload.ID,
		Name:       payload.Name,
		ProvinceID: payload.ProvinceID,
	}

	if err := s.db.Create(&region).Error; err != nil {
		return nil, err
	}

	return &region, nil
}

func (s *Service) EditRegionById(id string, payload PutRegionPayload) (error) {
	result := s.db.Model(&Region{}).Where("id = ?", id).Updates(payload)

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

// SOFT DELETE
func (s *Service) DeleteRegionById(id string) error {
	result := s.db.Delete(&Region{}, id)

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
