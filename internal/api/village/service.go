package village

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/utils"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

// PAYLOAD FROM REQUEST BODY
type VillagePayload struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DistrictID string `json:"district_id"`
}

func (s *Service) GetVillages() ([]Village, error) {
	var villages []Village

	if err := s.db.
		Preload("District.Region.Province").
		Find(&villages).Error; err != nil {
		return nil, err
	}

	return villages, nil
}

func (s *Service) GetVillageById(id string) (*Village, error) {
	var village Village

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

func (s *Service) AddVillage(payload VillagePayload) (*Village, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"ID":         payload.ID,
		"DistrictID": payload.DistrictID,
		"Name":       payload.Name,
	})

	if err != nil {
		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var dist district.District

	if err := s.db.First(&dist, payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}
		return nil, err
	}

	village := Village{
		ID:         payload.ID,
		Name:       payload.Name,
		DistrictID: payload.DistrictID,
	}

	if err := s.db.Create(&village).Error; err != nil {
		return nil, err
	}

	return &village, nil
}

func (s *Service) EditVillageById(id string, payload VillagePayload) (*Village, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"DistrictID": payload.DistrictID,
		"Name":       payload.Name,
	})

	if err != nil {
		return nil, err
	}

	var village Village

	// ENSURE VILLAGE EXISTS
	if err := s.db.First(&village, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}

		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var dist district.District

	if err := s.db.First(&dist, payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}

		return nil, err
	}

	village.Name = payload.Name
	village.DistrictID = payload.DistrictID

	if err := s.db.Save(&village).Error; err != nil {
		return nil, err
	}

	return &village, nil
}

// SOFT DELETE
func (s *Service) DeleteVillageById(id string) error {
	var village Village

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
