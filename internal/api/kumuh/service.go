package kumuh

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
type KumuhPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetKumuh() ([]KumuhPayload, error) {
	var kumuh []KumuhPayload

	if err := s.db.Find(&kumuh).Error; err != nil {
		return nil, err
	}

	return kumuh, nil
}
