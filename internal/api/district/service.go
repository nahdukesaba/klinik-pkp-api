package district

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
type DistrictPayload struct {
	Name     string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
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

	newDistrict := District{
		Name:     payload.Name,
		RegionID: regionID,
	}

	if err := s.db.Create(&newDistrict).Error; err != nil {
		return nil, err
	}

	return &newDistrict, nil
}

func (s *Service) EditDistrictById(id uint64, payload DistrictPayload) error {
	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return err
	}

	newDistrict := District{
		Name:     payload.Name,
		RegionID: regionID,
	}

	result := s.db.Model(&District{}).Where("id = ?", id).Updates(&newDistrict)

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
func (s *Service) DeleteDistrictById(id uint64) error {
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
