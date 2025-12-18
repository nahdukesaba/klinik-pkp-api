package service

import (
	"errors"
	"klinik-api/models"

	"gorm.io/gorm"
)

type RegencyService struct {
	db *gorm.DB
}

func NewRegencyService(db *gorm.DB) *RegencyService {
	return &RegencyService{db: db}
}

// GetAll returns all regencies
func (s *RegencyService) GetAll(provinceID uint, includeProvince bool) ([]models.Regency, error) {
	var regencies []models.Regency
	query := s.db.Order("name ASC")

	if provinceID > 0 {
		query = query.Where("province_id = ?", provinceID)
	}

	if includeProvince {
		query = query.Preload("Province")
	}

	if err := query.Find(&regencies).Error; err != nil {
		return nil, err
	}

	return regencies, nil
}

// GetByID returns regency by ID
func (s *RegencyService) GetByID(id uint, includeProvince bool) (*models.Regency, error) {
	var regency models.Regency
	query := s.db

	if includeProvince {
		query = query.Preload("Province")
	}

	if err := query.First(&regency, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("regency not found")
		}
		return nil, err
	}

	return &regency, nil
}

// Create creates new regency (admin only)
func (s *RegencyService) Create(provinceID uint, name, code string) (*models.Regency, error) {
	if provinceID == 0 {
		return nil, errors.New("province_id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Verify province exists
	var province models.Province
	if err := s.db.First(&province, provinceID).Error; err != nil {
		return nil, errors.New("province not found")
	}

	regency := models.Regency{
		ProvinceID: provinceID,
		Name:       name,
		Code:       code,
	}

	if err := s.db.Create(&regency).Error; err != nil {
		return nil, err
	}

	// Load province
	s.db.Preload("Province").First(&regency, regency.ID)

	return &regency, nil
}

// Update updates regency (admin only)
func (s *RegencyService) Update(id uint, provinceID uint, name, code string) (*models.Regency, error) {
	var regency models.Regency
	if err := s.db.First(&regency, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("regency not found")
		}
		return nil, err
	}

	if provinceID > 0 {
		// Verify province exists
		var province models.Province
		if err := s.db.First(&province, provinceID).Error; err != nil {
			return nil, errors.New("province not found")
		}
		regency.ProvinceID = provinceID
	}
	if name != "" {
		regency.Name = name
	}
	if code != "" {
		regency.Code = code
	}

	if err := s.db.Save(&regency).Error; err != nil {
		return nil, err
	}

	// Load province
	s.db.Preload("Province").First(&regency, regency.ID)

	return &regency, nil
}

// Delete deletes regency (admin only)
func (s *RegencyService) Delete(id uint) error {
	var regency models.Regency
	if err := s.db.First(&regency, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("regency not found")
		}
		return err
	}

	if err := s.db.Delete(&regency).Error; err != nil {
		return err
	}

	return nil
}
