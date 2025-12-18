package service

import (
	"errors"
	"klinik-api/models"

	"gorm.io/gorm"
)

type ProvinceService struct {
	db *gorm.DB
}

func NewProvinceService(db *gorm.DB) *ProvinceService {
	return &ProvinceService{db: db}
}

// GetAll returns all provinces with optional regencies
func (s *ProvinceService) GetAll(includeRegencies bool) ([]models.Province, error) {
	var provinces []models.Province
	query := s.db.Order("name ASC")

	if includeRegencies {
		query = query.Preload("Regencies", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		})
	}

	if err := query.Find(&provinces).Error; err != nil {
		return nil, err
	}

	return provinces, nil
}

// GetByID returns province by ID
func (s *ProvinceService) GetByID(id uint, includeRegencies bool) (*models.Province, error) {
	var province models.Province
	query := s.db

	if includeRegencies {
		query = query.Preload("Regencies", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		})
	}

	if err := query.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("province not found")
		}
		return nil, err
	}

	return &province, nil
}

// Create creates new province (admin only)
func (s *ProvinceService) Create(name, code string) (*models.Province, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	province := models.Province{
		Name: name,
		Code: code,
	}

	if err := s.db.Create(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

// Update updates province (admin only)
func (s *ProvinceService) Update(id uint, name, code string) (*models.Province, error) {
	var province models.Province
	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("province not found")
		}
		return nil, err
	}

	if name != "" {
		province.Name = name
	}
	if code != "" {
		province.Code = code
	}

	if err := s.db.Save(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

// Delete deletes province (admin only)
func (s *ProvinceService) Delete(id uint) error {
	var province models.Province
	if err := s.db.First(&province, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("province not found")
		}
		return err
	}

	if err := s.db.Delete(&province).Error; err != nil {
		return err
	}

	return nil
}
