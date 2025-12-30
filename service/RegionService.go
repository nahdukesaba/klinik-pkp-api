package service

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"
)

// REGION SERVICE STRUCT
type RegionService struct {
	db        *gorm.DB
	validator *utils.Validator
}

// REPRESENTS DATA FROM INPUT
type RegionPayload struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProvinceID string `json:"province_id"`
}

// REGION SERVICE CONSTRUCTOR
func NewRegionService(db *gorm.DB) *RegionService {
	return &RegionService{db: db, validator: &utils.Validator{}}
}

func (s *RegionService) GetAll() ([]models.Region, error) {
	var regions []models.Region

	if err := s.db.Preload("Province").Find(&regions).Error; err != nil {
		return nil, err
	}

	return regions, nil
}

func (s *RegionService) GetByID(id string) (*models.Region, error) {
	var region models.Region

	if err := s.db.Preload("Province").First(&region, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}
		return nil, err
	}

	return &region, nil
}

func (s *RegionService) Create(input RegionPayload) (*models.Region, error) {
	if ok, msg := s.validator.ValidateRequired(input.ID, "ID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.ProvinceID, "ProvinceID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	// ENSURE RELATED PROVINCE EXISTS
	var province models.Province

	if err := s.db.First(&province, input.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	region := models.Region{
		ID:         input.ID,
		Name:       input.Name,
		ProvinceID: input.ProvinceID,
	}

	if err := s.db.Create(&region).Error; err != nil {
		return nil, err
	}

	return &region, nil
}

func (s *RegionService) Update(id string, input RegionPayload) (*models.Region, error) {
	if ok, msg := s.validator.ValidateRequired(input.ProvinceID, "ProvinceID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	var region models.Region

	// ENSURE REGION EXISTS
	if err := s.db.First(&region, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var province models.Province

	if err := s.db.First(&province, input.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	region.Name = input.Name
	region.ProvinceID = input.ProvinceID

	if err := s.db.Save(&region).Error; err != nil {
		return nil, err
	}

	return &region, nil
}

// SOFT DELETE
func (s *RegionService) Delete(id string) error {
	var region models.Region

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
