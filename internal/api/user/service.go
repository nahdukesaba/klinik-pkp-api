package user

import (
	"errors"
	"klinik-pkp-api/utils"

	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	validator *utils.Validator
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:        db,
		validator: &utils.Validator{},
	}
}

// REGISTER payload STRUCT
type RegisterPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"` // Optional
}

// LOGIN payload STRUCT
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// CREATE NEW USER
func (s *Service) Register(payload RegisterPayload) (*AuthResponse, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"Name":     payload.Name,
		"Email":    payload.Email,
		"Password": payload.Password,
	})

	if err != nil {
		return nil, err
	}

	// VALIDATE EMAIL FORMAT
	if !s.validator.ValidateEmail(payload.Email) {
		return nil, errors.New("invalid email format")
	}

	// VALIDATE PHONE FORMAT IF PROVIDED
	if payload.Phone != "" && !s.validator.ValidatePhone(payload.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// VALIDATE PASSWORD STRENGTH (MIN 8 CHARS)
	if ok, message := s.validator.ValidatePasswordStrength(payload.Password); !ok {
		return nil, errors.New(message)
	}

	// CHECK IF EMAIL ALREADY EXISTS
	var existingUser User

	if err := s.db.Where("email = ?", payload.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// CHECK IF PHONE ALREADY EXISTS
	if payload.Phone != "" {
		if err := s.db.Where("phone = ?", payload.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("phone number already registered")
		}
	}

	// CREATE USER RECORD IN DB
	user := User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Phone:    payload.Phone,
		Role:     "User",
		IsActive: true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// GENERATE TOKEN
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  &user,
	}, nil
}

// LOGIN USER
func (s *Service) Login(payload LoginPayload) (*AuthResponse, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"Email":    payload.Email,
		"Password": payload.Password,
	})

	if err != nil {
		return nil, err
	}

	// Find user by email
	var user User
	if err := s.db.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(payload.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  &user,
	}, nil
}

// RETURN ALL USERS
func (s *Service) GetAll() ([]User, error) {
	var users []User

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GET PROFILE BY ID
func (s *Service) GetProfile(userID uint) (*User, error) {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &user, nil
}

// UPDATE PROFILE BY ID
func (s *Service) UpdateProfile(userID uint, name string) (*User, error) {
	// VALIDATIONS
	err := s.validator.ValidateRequiredFields(map[string]any{
		"Name": name,
	})

	if err != nil {
		return nil, err
	}

	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = name
	user.UpdatedBy = &userID
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
