package district

import (
	"errors"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

// Service struct
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// DistrictPayload represents data from input
type DistrictPayload struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RegionID string `json:"region_id"`
}

// NewService constructor
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetAll() ([]District, error) {
	var districts []District

	if err := s.db.Preload("Region.Province").Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

func (s *Service) GetByID(id string) (*District, error) {
	var district District

	if err := s.db.Preload("Region.Province").First(&district, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}

		return nil, err
	}

	return &district, nil
}

func (s *Service) Create(payload DistrictPayload) (*District, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"ID":       payload.ID,
		"RegionID": payload.RegionID,
		"Name":     payload.Name,
	})

	if err != nil {
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var reg region.Region

	if err := s.db.First(&reg, payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

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

func (s *Service) Update(id string, payload DistrictPayload) (*District, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"RegionID": payload.RegionID,
		"Name":     payload.Name,
	})

	if err != nil {
		return nil, err
	}

	var district District

	// ENSURE DISTRICT EXISTS
	if err := s.db.First(&district, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}

		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var reg region.Region

	if err := s.db.First(&reg, payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	district.Name = payload.Name
	district.RegionID = payload.RegionID

	if err := s.db.Save(&district).Error; err != nil {
		return nil, err
	}

	return &district, nil
}

// SOFT DELETE
func (s *Service) Delete(id string) error {
	var district District

	if err := s.db.First(&district, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("District not found")
		}

		return err
	}

	if err := s.db.Delete(&district).Error; err != nil {
		return err
	}

	return nil
}
