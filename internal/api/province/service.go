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

// POST PROVINCE PAYLOAD FROM REQUEST BODY
type PostProvincePayload struct {
	ID   string `json:"id" validate:"required,min=2,max=3"`
	Name string `json:"name" validate:"required"`
}

// PUT PROVINCE PAYLOAD FROM REQUEST BODY
type PutProvincePayload struct {
	Name string `json:"name" validate:"required"`
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

func (s *Service) GetProvinceById(id string) (*Province, error) {
	var province Province

	if err := s.db.First(&province, id).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) AddProvince(payload PostProvincePayload) (*Province, error) {
	province := Province{
		ID:   payload.ID,
		Name: payload.Name,
	}

	if err := s.db.Create(&province).Error; err != nil {
		return nil, err
	}

	return &province, nil
}

func (s *Service) EditProvinceById(id string, payload PutProvincePayload) error {
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

func (s *Service) DeleteProvinceById(id string) error {
	result := s.db.Delete(&Province{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
