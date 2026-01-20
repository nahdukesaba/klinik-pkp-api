package village

import (
	"gorm.io/gorm"
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

// POST VILLAGE PAYLOAD FROM REQUEST BODY
type PostVillagePayload struct {
	ID         string `json:"id" validate:"required,min=10,max=10"`
	Name       string `json:"name" validate:"required"`
	DistrictID string `json:"district_id" validate:"required"`
}

// PUT VILLAGE PAYLOAD FROM REQUEST BODY
type PutVillagePayload struct {
	Name       string `json:"name" validate:"required"`
	DistrictID string `json:"district_id" validate:"required"`
}

func (s *Service) GetVillages() ([]Village, error) {
	var villages []Village

	if err := s.db.Preload("District.Region.Province").Find(&villages).Error; err != nil {
		return nil, err
	}

	return villages, nil
}

func (s *Service) GetVillageById(id string) (*Village, error) {
	var village Village

	if err := s.db.Preload("District.Region.Province").First(&village, id).Error; err != nil {
		return nil, err
	}

	return &village, nil
}

func (s *Service) AddVillage(payload PostVillagePayload) (*Village, error) {
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

func (s *Service) EditVillageById(id string, payload PutVillagePayload) error {
	result := s.db.Model(&Village{}).Where("id = ?", id).Updates(payload)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// SOFT DELETE
func (s *Service) DeleteVillageById(id string) error {
	result := s.db.Delete(&Village{}, id)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
