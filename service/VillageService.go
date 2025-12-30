package service

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"
)

// VILLAGE SERVICE STRUCT
type VillageService struct {
	db        *gorm.DB
	validator *utils.Validator
}

// VILLAGE SERVICE CONSTRUCTOR
func NewVillageService(db *gorm.DB) *VillageService {
	return &VillageService{db: db, validator: &utils.Validator{}}
}

// VILLAGE INPUT STRUCT
type VillagePayload struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DistrictID string `json:"district_id"`
}

func (s *VillageService) GetAll() ([]models.Village, error) {
	var villages []models.Village

	if err := s.db.
		Preload("District.Region.Province").
		Find(&villages).Error; err != nil {
		return nil, err
	}

	return villages, nil
}

func (s *VillageService) GetByID(id string) (*models.Village, error) {
	var village models.Village

	if err := s.db.
		Preload("District.Region.Province").
		First(&village, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}

		return nil, err
	}

	return &village, nil
}

func (s *VillageService) Create(input VillagePayload) (*models.Village, error) {
	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.DistrictID, "DistrictID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.ID, "ID"); !ok {
		return nil, errors.New(msg)
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district models.District

	if err := s.db.First(&district, input.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}
		return nil, err
	}

	village := models.Village{
		ID:         input.ID,
		Name:       input.Name,
		DistrictID: input.DistrictID,
	}

	if err := s.db.Create(&village).Error; err != nil {
		return nil, err
	}

	return &village, nil
}

func (s *VillageService) Update(id string, input VillagePayload) (*models.Village, error) {
	if ok, msg := s.validator.ValidateRequired(input.DistrictID, "DistrictID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	var village models.Village

	// ENSURE VILLAGE EXISTS
	if err := s.db.First(&village, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}

		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district models.District

	if err := s.db.First(&district, input.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}

		return nil, err
	}

	village.Name = input.Name
	village.DistrictID = input.DistrictID

	if err := s.db.Save(&village).Error; err != nil {
		return nil, err
	}

	return &village, nil
}

// SOFT DELETE
func (s *VillageService) Delete(id string) error {
	var village models.Village

	if err := s.db.First(&village, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Village not found")
		}

		return err
	}

	if err := s.db.Delete(&village).Error; err != nil {
		return err
	}

	return nil
}
