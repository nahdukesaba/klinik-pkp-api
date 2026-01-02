package region

import (
	"errors"
	"klinik-pkp-api/internal/api/province"
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

// Service struct
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// RegionPayload represents data from input
type RegionPayload struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProvinceID string `json:"province_id"`
}

// NewService constructor
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetAll() ([]Region, error) {
	var regions []Region

	if err := s.db.Preload("Province").Find(&regions).Error; err != nil {
		return nil, err
	}

	return regions, nil
}

func (s *Service) GetByID(id string) (*Region, error) {
	var region Region

	if err := s.db.Preload("Province").First(&region, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}
		return nil, err
	}

	return &region, nil
}

func (s *Service) Create(payload RegionPayload) (*Region, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"ID":         payload.ID,
		"ProvinceID": payload.ProvinceID,
		"Name":       payload.Name,
	})

	if err != nil {
		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var prov province.Province

	if err := s.db.First(&prov, payload.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

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

func (s *Service) Update(id string, payload RegionPayload) (*Region, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"ProvinceID": payload.ProvinceID,
		"Name":       payload.Name,
	})

	if err != nil {
		return nil, err
	}

	var region Region

	// ENSURE REGION EXISTS
	if err := s.db.First(&region, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var prov province.Province

	if err := s.db.First(&prov, payload.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	region.Name = payload.Name
	region.ProvinceID = payload.ProvinceID

	if err := s.db.Save(&region).Error; err != nil {
		return nil, err
	}

	return &region, nil
}

// SOFT DELETE
func (s *Service) Delete(id string) error {
	var region Region

	if err := s.db.First(&region, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Region not found")
		}

		return err
	}

	if err := s.db.Delete(&region).Error; err != nil {
		return err
	}

	return nil
}
