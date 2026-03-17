package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FOR UPDATING USER PROFILE
type EditUserProfilePayload struct {
	Name  string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
	Phone string `json:"phone" validate:"-"`
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		validator: utils.NewValidator(),
	}
}

func (s *Service) GetUsers() ([]User, error) {
	var users []User

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) GetUserById(id uuid.UUID) (*User, error) {
	var user User

	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) EditUserById(id uuid.UUID, payload *EditUserProfilePayload) error {
	// CHECK IF USER EXISTS
	existingUser := User{}

	if err := s.db.First(&existingUser, id).Error; err != nil {
		return err
	}

	newUser := User{
		Name:  payload.Name,
		Phone: payload.Phone,
	}

	result := s.db.Model(&User{}).Where("id = ?", id).Updates(newUser)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
