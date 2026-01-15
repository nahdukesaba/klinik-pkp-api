package province

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type ProvincePayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetProvinces() ([]Province, error) {
	var provinces []Province

	if err := s.db.Find(&provinces).Error; err != nil {
		return nil, err
	}

	return provinces, nil
}

func (s *Service) GetProvinceById(id string) (*Province, error) {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("province with id %s is not found", id)
		}
		return nil, err
	}

	return &province, nil
}

func (s *Service) AddProvince(payload ProvincePayload) (*Province, error) {
	// VALIDATIONS
	if err := s.validator.ValidateRequiredFields(map[string]any{
		"ID":   payload.ID,
		"Name": payload.Name,
	}); err != nil {
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

func (s *Service) EditProvinceById(id string, payload ProvincePayload) (*Province, error) {
	// VALIDATIONS
	if err := s.validator.ValidateRequiredFields(map[string]any{
		"Name": payload.Name,
	}); err != nil {
		return nil, err
	}

	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("province with id %s is not found", id)
		}

		return nil, err
	}

	province.Name = payload.Name

	if err := s.db.Save(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) DeleteProvinceById(id string) error {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("province with id %s is not found", id)
		}
		return err
	}

	if err := s.db.Delete(&province).Error; err != nil {
		return err
	}

	return nil
}
