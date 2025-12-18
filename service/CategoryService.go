package service

import (
	"errors"
	"klinik-api/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

// GetAll returns all categories
func (s *CategoryService) GetAll() ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetByID returns category by ID
func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// Create creates new category (admin only)
func (s *CategoryService) Create(name, description string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	category := models.Category{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// Update updates category (admin only)
func (s *CategoryService) Update(id uint, name, description string) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	if name != "" {
		category.Name = name
	}
	if description != "" {
		category.Description = description
	}

	if err := s.db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

// Delete deletes category (admin only)
func (s *CategoryService) Delete(id uint) error {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("category not found")
		}
		return err
	}

	if err := s.db.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
