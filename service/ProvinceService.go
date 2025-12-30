package service

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"
)

// PROVINCE SERVICE STRUCT
type ProvinceService struct {
	db        *gorm.DB
	validator *utils.Validator
}

// REPRESENTS DATA FROM INPUT
type ProvincePayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PROVINCE SERVICE CONSTRUCTOR
func NewProvinceService(db *gorm.DB) *ProvinceService {
	return &ProvinceService{db: db, validator: &utils.Validator{}}
}

func (s *ProvinceService) GetAll() ([]models.Province, error) {
	var provinces []models.Province

	if err := s.db.Find(&provinces).Error; err != nil {
		return nil, err
	}

	return provinces, nil
}

func (s *ProvinceService) GetByID(id string) (*models.Province, error) {
	var province models.Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}
		return nil, err
	}

	return &province, nil
}

func (s *ProvinceService) Create(input ProvincePayload) (*models.Province, error) {
	if ok, msg := s.validator.ValidateRequired(input.ID, "ID"); !ok {
		return nil, errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	province := models.Province{
		ID:   input.ID,
		Name: input.Name,
	}

	if err := s.db.Create(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *ProvinceService) Update(id string, input ProvincePayload) (*models.Province, error) {
	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	var province models.Province

	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	province.Name = input.Name

	if err := s.db.Save(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *ProvinceService) Delete(id string) error {
	var province models.Province

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
