package service

import (
	"errors"
	"klinik-api/config"
	"klinik-api/models"
	"klinik-api/utils"

	"gorm.io/gorm"
)

type AdminService struct {
	db        *gorm.DB
	validator *utils.Validator
	config    *config.Config
}

func NewAdminService(db *gorm.DB, cfg *config.Config) *AdminService {
	return &AdminService{
		db:        db,
		validator: &utils.Validator{},
		config:    cfg,
	}
}

// CreateUserInput for admin to create users
type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"` // Optional
	Role     string `json:"role"`
}

// UpdateUserInput for admin to update users
type UpdateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"` // Optional
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

// AdminLoginInput represents admin login request
type AdminLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AdminAuthResponse represents authentication response
type AdminAuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login authenticates admin and returns token
func (s *AdminService) Login(input AdminLoginInput) (*AdminAuthResponse, error) {
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

	// Check if user has admin role
	if !user.IsAdmin() {
		return nil, errors.New("access denied: admin role required")
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

	return &AdminAuthResponse{
		Token: token,
		User:  &user,
	}, nil
}

// GetAllUsers returns all users with pagination
func (s *AdminService) GetAllUsers(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID returns user by ID
func (s *AdminService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates new user (admin only)
func (s *AdminService) CreateUser(input CreateUserInput, createdBy uint) (*models.User, error) {
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
	if ok, msg := s.validator.ValidateRequired(input.Role, "Role"); !ok {
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

	// Validate role - Public is not a valid database role (it's for non-authenticated users)
	validRoles := []string{
		s.config.RoleSuperAdmin,
		s.config.RoleAdminEselon1,
		s.config.RoleVerificatorEselon1,
		s.config.RoleAdminBalai,
		s.config.RoleVerificatorBalai,
		s.config.RoleSurveyor,
		s.config.RoleUser,
	}
	
	isValidRole := false
	for _, role := range validRoles {
		if input.Role == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return nil, errors.New("invalid role")
	}

	// Create user with audit fields
	user := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password, // Will be hashed by BeforeCreate hook
		Phone:     input.Phone,
		Role:      input.Role,
		IsActive:  true,
		CreatedBy: &createdBy,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates user (admin only)
func (s *AdminService) UpdateUser(id uint, input UpdateUserInput, updatedBy uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		if !s.validator.ValidateEmail(input.Email) {
			return nil, errors.New("invalid email format")
		}
		// Check if email already exists
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", input.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already registered")
		}
		user.Email = input.Email
	}
	if input.Phone != "" {
		if !s.validator.ValidatePhone(input.Phone) {
			return nil, errors.New("invalid phone number format")
		}
		// Check if phone already exists
		var existingUser models.User
		if err := s.db.Where("phone = ? AND id != ?", input.Phone, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("phone number already registered")
		}
		user.Phone = input.Phone
	}
	if input.Role != "" {
		user.Role = input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	
	user.UpdatedBy = &updatedBy

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes user (admin only)
func (s *AdminService) DeleteUser(id uint, deletedBy uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}
	
	// Set who deleted this record
	user.DeletedBy = &deletedBy
	s.db.Save(&user)

	// Soft delete
	if err := s.db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
