package user

import (
	"klinik-pkp-api/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

// PAYLOAD FOR CREATING NEW USER (ADMIN ONLY)
type AddUserPayload struct {
	Name     string `json:"name" validate:"required,ne=null,ne=NULL,ne=Null"`
	Email    string `json:"email" validate:"required,email"`
	NIP      string `json:"nip" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"-"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	IsActive bool   `json:"is_active" validate:"required"`
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

func (s *Service) GetUsers(pagination utils.PaginationQuery) ([]User, int64, error) {
	var users []User
	var totalRecords int64

	query := s.db.Model(&User{})

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(pagination.Limit).Offset(pagination.Offset()).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalRecords, nil
}

func (s *Service) GetUserById(id uuid.UUID) (*User, error) {
	var user User

	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) AddUser(payload *AddUserPayload) (*User, error) {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	newUser := User{
		ID:       uuid.New(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		NIP:      payload.NIP,
		Phone:    payload.Phone,
		Role:     payload.Role,
		IsActive: payload.IsActive,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *Service) EditUserById(id uuid.UUID, payload *EditUserProfilePayload) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// CHECK IF USER EXISTS
	existingUser := User{}

	if err := tx.First(&existingUser, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	newUser := User{
		Name:  payload.Name,
		Phone: payload.Phone,
	}

	result := tx.Model(&User{}).Where("id = ?", id).Updates(newUser)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
