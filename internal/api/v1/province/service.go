package province

import (
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FROM REQUEST BODY
type ProvincePayload struct {
	Name string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: utils.NewValidator()}
}

func (s *Service) GetProvinces() ([]Province, error) {
	var provinces []Province

	if err := s.db.Find(&provinces).Error; err != nil {
		return nil, err
	}

	return provinces, nil
}

func (s *Service) GetProvinceById(id uint64) (*Province, error) {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) AddProvince(payload ProvincePayload) (*Province, error) {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	newProvince := Province{
		Name: payload.Name,
	}

	if err := tx.Create(&newProvince).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newProvince, nil
}

func (s *Service) EditProvinceById(id uint64, payload ProvincePayload) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// USE .Update() FOR SINGLE FIELD UPDATE
	// AVOID .Save() TO PREVENT OVERWRITING OTHER FIELDS WITH ZERO VALUES
	result := tx.Model(&Province{}).Where("id = ?", id).Update("name", payload.Name)

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

func (s *Service) DeleteProvinceById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&Province{}, id)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

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
