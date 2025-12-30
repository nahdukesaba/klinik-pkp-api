package service

import (
	"errors"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

// DistrictService handles district CRUD operations
type DistrictService struct {
	db        *gorm.DB
	validator *utils.Validator
}

// REPRESENTS DATA FROM INPUT
type DistrictPayload struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RegionID string `json:"region_id"`
}

// DISTRICT SERVICE CONTROLLER
func NewDistrictService(db *gorm.DB) *DistrictService {
	return &DistrictService{db: db, validator: &utils.Validator{}}
}

func (s *DistrictService) GetAll() ([]models.District, error) {
	var districts []models.District

	if err := s.db.Preload("Region.Province").Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

func (s *DistrictService) GetByID(id string) (*models.District, error) {
	var district models.District

	if err := s.db.Preload("Region.Province").First(&district, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}

		return nil, err
	}

	return &district, nil
}

func (s *DistrictService) Create(input DistrictPayload) (*models.District, error) {
	if ok, msg := s.validator.ValidateRequired(input.ID, "ID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.RegionID, "RegionID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	// ENSURE RELATED REGION EXISTS
	var region models.Region

	if err := s.db.First(&region, input.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	district := models.District{
		ID:       input.ID,
		Name:     input.Name,
		RegionID: input.RegionID,
	}

	if err := s.db.Create(&district).Error; err != nil {
		return nil, err
	}

	return &district, nil
}

func (s *DistrictService) Update(id string, input DistrictPayload) (*models.District, error) {
	if ok, msg := s.validator.ValidateRequired(input.RegionID, "RegionID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	var district models.District

	// ENSURE DISTRICT EXISTS
	if err := s.db.First(&district, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}

		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region models.Region

	if err := s.db.First(&region, input.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	district.Name = input.Name
	district.RegionID = input.RegionID

	if err := s.db.Save(&district).Error; err != nil {
		return nil, err
	}

	return &district, nil
}

// SOFT DELETE
func (s *DistrictService) Delete(id string) error {
	var district models.District

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
