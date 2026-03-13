package village

import (
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
	"strconv"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type VillagePayload struct {
	Name       string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
	DistrictID string `json:"district_id" validate:"required"`
}

// FILTER FOR QUERY PARAMETERS
type VillageFilter struct {
	DistrictID *uint64
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: utils.NewValidator()}
}

func (s *Service) GetVillages(filter VillageFilter) ([]Village, error) {
	var villages []Village

	query := s.db.Preload("District.Region.Province").Model(&Village{})

	// FILTER BY DISTRICT
	if filter.DistrictID != nil {
		query = query.Where("village.district_id = ?", *filter.DistrictID)
	}

	if err := query.Find(&villages).Error; err != nil {
		return nil, err
	}

	return villages, nil
}

func (s *Service) GetVillageById(id uint64) (*Village, error) {
	var village Village

	if err := s.db.Preload("District.Region.Province").First(&village, id).Error; err != nil {
		return nil, err
	}

	return &village, nil
}

func (s *Service) AddVillage(payload VillagePayload) (*Village, error) {
	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return nil, err
	}

	newVillage := Village{
		Name:       payload.Name,
		DistrictID: districtID,
	}

	if err := s.db.Create(&newVillage).Error; err != nil {
		return nil, err
	}

	return &newVillage, nil
}

func (s *Service) EditVillageById(id uint64, payload VillagePayload) error {
	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return err
	}

	newVillage := Village{
		Name:       payload.Name,
		DistrictID: districtID,
	}

	result := s.db.Model(&Village{}).Where("id = ?", id).Updates(&newVillage)

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
func (s *Service) DeleteVillageById(id uint64) error {
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
