package sosialisasi

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
type SosialisasiPayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Location    string `json:"location"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB) *Service {
	return &Service{db: db, validator: &utils.Validator{}}
}

func (s *Service) GetSosialisasi() ([]Sosialisasi, error) {
	var sosialisasi []Sosialisasi

	if err := s.db.Find(&sosialisasi).Error; err != nil {
		return nil, err
	}

	return sosialisasi, nil
}
