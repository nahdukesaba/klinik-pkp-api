package service

import (
	"errors"
	"klinik-api/models"
	"klinik-api/utils"

	"gorm.io/gorm"
)

type UserService struct {
	db        *gorm.DB
	validator *utils.Validator
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:        db,
		validator: &utils.Validator{},
	}
}

// RegisterInput represents user registration request
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"` // Optional
}

// LoginInput represents user login request
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Register creates new user account (User role by default)
func (s *UserService) Register(input RegisterInput) (*AuthResponse, error) {
	// Validate required fields
	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.Email, "Email"); !ok {
		return nil, errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.Password, "Password"); !ok {
		return nil, errors.New(msg)
	}

	// Validate email format
	if !s.validator.ValidateEmail(input.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate phone format if provided
	if input.Phone != "" && !s.validator.ValidatePhone(input.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// Validate password strength (min 8 chars)
	if ok, msg := s.validator.ValidatePasswordStrength(input.Password); !ok {
		return nil, errors.New(msg)
	}

	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Check if phone already exists (if provided)
	if input.Phone != "" {
		if err := s.db.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("phone number already registered")
		}
	}

	// Create user with User role
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password, // Will be hashed by BeforeCreate hook
		Phone:    input.Phone,
		Role:     "User",
		IsActive: true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
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

// Login authenticates user and returns token
func (s *UserService) Login(input LoginInput) (*AuthResponse, error) {
	// Validate required fields
	if ok, msg := s.validator.ValidateRequired(input.Email, "Email"); !ok {
		return nil, errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.Password, "Password"); !ok {
		return nil, errors.New(msg)
	}

	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
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
	if !user.CheckPassword(input.Password) {
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

// GetProfile returns user profile by ID
func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uint, name string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if ok, msg := s.validator.ValidateRequired(name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	user.Name = name
	user.UpdatedBy = &userID
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ChangePasswordInput represents password change request
type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword updates user password with old password validation
func (s *UserService) ChangePassword(userID uint, input ChangePasswordInput) error {
	// Validate required fields
	if ok, msg := s.validator.ValidateRequired(input.OldPassword, "Old Password"); !ok {
		return errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.NewPassword, "New Password"); !ok {
		return errors.New(msg)
	}

	// Validate new password strength
	if ok, msg := s.validator.ValidatePasswordStrength(input.NewPassword); !ok {
		return errors.New(msg)
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !user.CheckPassword(input.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := user.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedBy = &userID
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
