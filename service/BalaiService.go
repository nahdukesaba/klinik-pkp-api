package service

import (
	"errors"
	"klinik-api/models"

	"gorm.io/gorm"
)

type BalaiService struct {
	db *gorm.DB
}

func NewBalaiService(db *gorm.DB) *BalaiService {
	return &BalaiService{db: db}
}

// GetAll returns all balais
func (s *BalaiService) GetAll(provinceID, categoryID uint) ([]models.Balai, error) {
	var balais []models.Balai
	query := s.db.Preload("Province").Preload("Category").Order("name ASC")

	if provinceID > 0 {
		query = query.Where("province_id = ?", provinceID)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Find(&balais).Error; err != nil {
		return nil, err
	}

	return balais, nil
}

// GetByID returns balai by ID
func (s *BalaiService) GetByID(id uint) (*models.Balai, error) {
	var balai models.Balai
	if err := s.db.Preload("Province").Preload("Category").First(&balai, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("balai not found")
		}
		return nil, err
	}
	return &balai, nil
}

// Create creates new balai (admin only)
func (s *BalaiService) Create(provinceID, categoryID uint, name, address string) (*models.Balai, error) {
	if provinceID == 0 {
		return nil, errors.New("province_id is required")
	}
	if categoryID == 0 {
		return nil, errors.New("category_id is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Verify province exists
	var province models.Province
	if err := s.db.First(&province, provinceID).Error; err != nil {
		return nil, errors.New("province not found")
	}

	// Verify category exists
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		return nil, errors.New("category not found")
	}

	balai := models.Balai{
		ProvinceID: provinceID,
		CategoryID: categoryID,
		Name:       name,
		Address:    address,
	}

	if err := s.db.Create(&balai).Error; err != nil {
		return nil, err
	}

	// Load relationships
	s.db.Preload("Province").Preload("Category").First(&balai, balai.ID)

	return &balai, nil
}

// Update updates balai (admin only)
func (s *BalaiService) Update(id, provinceID, categoryID uint, name, address string) (*models.Balai, error) {
	var balai models.Balai
	if err := s.db.First(&balai, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("balai not found")
		}
		return nil, err
	}

	if provinceID > 0 {
		var province models.Province
		if err := s.db.First(&province, provinceID).Error; err != nil {
			return nil, errors.New("province not found")
		}
		balai.ProvinceID = provinceID
	}
	if categoryID > 0 {
		var category models.Category
		if err := s.db.First(&category, categoryID).Error; err != nil {
			return nil, errors.New("category not found")
		}
		balai.CategoryID = categoryID
	}
	if name != "" {
		balai.Name = name
	}
	if address != "" {
		balai.Address = address
	}

	if err := s.db.Save(&balai).Error; err != nil {
		return nil, err
	}

	// Load relationships
	s.db.Preload("Province").Preload("Category").First(&balai, balai.ID)

	return &balai, nil
}

// Delete deletes balai (admin only)
func (s *BalaiService) Delete(id uint) error {
	var balai models.Balai
	if err := s.db.First(&balai, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("balai not found")
		}
		return err
	}

	if err := s.db.Delete(&balai).Error; err != nil {
		return err
	}

	return nil
}
