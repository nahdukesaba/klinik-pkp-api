package province

import (
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
	Name string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
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

func (s *Service) GetProvinceById(id uint64) (*Province, error) {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) AddProvince(payload ProvincePayload) (*Province, error) {
	newProvince := Province{
		Name: payload.Name,
	}

	if err := s.db.Create(&newProvince).Error; err != nil {
		return nil, err
	}

	return &newProvince, nil
}

func (s *Service) EditProvinceById(id uint64, payload ProvincePayload) error {
	// USE .Update() FOR SINGLE FIELD UPDATE
	// AVOID .Save() TO PREVENT OVERWRITING OTHER FIELDS WITH ZERO VALUES
	result := s.db.Model(&Province{}).Where("id = ?", id).Update("name", payload.Name)

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

func (s *Service) DeleteProvinceById(id uint64) error {
	result := s.db.Delete(&Province{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
