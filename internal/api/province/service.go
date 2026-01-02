package province

import (
	"errors"
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

// Service struct
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// ProvincePayload represents data from input
type ProvincePayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewService constructor
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetAll() ([]Province, error) {
	var provinces []Province

	if err := s.db.Find(&provinces).Error; err != nil {
		return nil, err
	}

	return provinces, nil
}

func (s *Service) GetByID(id string) (*Province, error) {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}
		return nil, err
	}

	return &province, nil
}

func (s *Service) Create(payload ProvincePayload) (*Province, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"ID":   payload.ID,
		"Name": payload.Name,
	})

	if err != nil {
		return nil, err
	}

	province := Province{
		ID:   payload.ID,
		Name: payload.Name,
	}

	if err := s.db.Create(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) Update(id string, payload ProvincePayload) (*Province, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"Name": payload.Name,
	})

	if err != nil {
		return nil, err
	}

	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	province.Name = payload.Name

	if err := s.db.Save(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) Delete(id string) error {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Province not found")
		}
		return err
	}

	if err := s.db.Delete(&province).Error; err != nil {
		return err
	}

	return nil
}
