package service

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"
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

// REGISTER INPUT STRUCT
type RegisterPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"` // Optional
}

// LOGIN INPUT STRUCT
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// CREATE NEW USER
func (s *UserService) Register(input RegisterPayload) (*AuthResponse, error) {

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.Email, "Email"); !ok {
		return nil, errors.New(msg)
	}
	if ok, msg := s.validator.ValidateRequired(input.Password, "Password"); !ok {
		return nil, errors.New(msg)
	}

	// VALIDATE EMAIL FORMAT
	if !s.validator.ValidateEmail(input.Email) {
		return nil, errors.New("invalid email format")
	}

	// VALIDATE PHONE FORMAT IF PROVIDED
	if input.Phone != "" && !s.validator.ValidatePhone(input.Phone) {
		return nil, errors.New("invalid phone number format")
	}

	// VALIDATE PASSWORD STRENGTH (MIN 8 CHARS)
	if ok, msg := s.validator.ValidatePasswordStrength(input.Password); !ok {
		return nil, errors.New(msg)
	}

	// CHECK IF EMAIL ALREADY EXISTS
	var existingUser models.User

	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// CHECK IF PHONE ALREADY EXISTS
	if input.Phone != "" {
		if err := s.db.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("phone number already registered")
		}
	}

	// CREATE USER RECORD IN DB
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Phone:    input.Phone,
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
func (s *UserService) Login(input LoginPayload) (*AuthResponse, error) {
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

// RETURN ALL USERS
func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil

	// users := make([]models.User, 0)

	// err := s.db.
	// 	Model(&models.User{}).
	// 	Find(&users).
	// 	Error

	// if err != nil {
	// 	return nil, err
	// }

	// return users, nil
}

// GET PROFILE BY ID
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

// UPDATE PROFILE BY ID
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

// CHANGE PASSWORD INPUT STRUCT
type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// CHANGE PASSWORD
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
