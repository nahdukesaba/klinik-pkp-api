package district

import (
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// POST DISTRICT PAYLOAD FROM REQUEST BODY
type PostDistrictPayload struct {
	ID       string `json:"id" validate:"required,min=6,max=6"`
	Name     string `json:"name" validate:"required"`
	RegionID string `json:"region_id" validate:"required"`
}

// PUT DISTRICT PAYLOAD FROM REQUEST BODY
type PutDistrictPayload struct {
	Name     string `json:"name" validate:"required"`
	RegionID string `json:"region_id" validate:"required"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetDistricts() ([]District, error) {
	var districts []District

	if err := s.db.Preload("Region.Province").Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

func (s *Service) GetDistrictById(id string) (*District, error) {
	var district District

	if err := s.db.Preload("Region.Province").First(&district, id).Error; err != nil {
		return nil, err
	}

	return &district, nil
}

func (s *Service) AddDistrict(payload PostDistrictPayload) (*District, error) {
	district := District{
		ID:       payload.ID,
		Name:     payload.Name,
		RegionID: payload.RegionID,
	}

	if err := s.db.Create(&district).Error; err != nil {
		return nil, err
	}

	return &district, nil
}

func (s *Service) EditDistrictById(id string, payload PutDistrictPayload) error {
	result := s.db.Model(&District{}).Where("id = ?", id).Updates(payload)

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
func (s *Service) DeleteDistrictById(id string) error {
	result := s.db.Delete(&District{}, id)

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
