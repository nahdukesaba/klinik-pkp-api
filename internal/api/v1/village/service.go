package village

import (
	"klinik-pkp-api/utils"
	"strconv"

	"gorm.io/gorm"
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

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	newVillage := Village{
		Name:       payload.Name,
		DistrictID: districtID,
	}

	if err := tx.Create(&newVillage).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newVillage, nil
}

func (s *Service) EditVillageById(id uint64, payload VillagePayload) error {
	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	newVillage := Village{
		Name:       payload.Name,
		DistrictID: districtID,
	}

	result := tx.Model(&Village{}).Where("id = ?", id).Updates(&newVillage)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// SOFT DELETE
func (s *Service) DeleteVillageById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&Village{}, id)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
