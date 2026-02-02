package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents user account in the system
type User struct {
	ID       uuid.UUID `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"size:255;not null" json:"name"`
	Email    string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string    `gorm:"size:255;not null" json:"-"`                 // HIDE FROM JSON
	Phone    string    `gorm:"size:20;uniqueIndex" json:"phone,omitempty"` // UNIQUE
	Role     string    `gorm:"type:varchar(50);not null" json:"role"`
	IsActive bool      `gorm:"default:true" json:"is_active"`

	CreatedBy *uint `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy *uint `gorm:"index" json:"deleted_by,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword verifies password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HashPassword hashes a password
func (u *User) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsSuperAdmin checks if user is super admin (alias untuk IsAdmin)
func (u *User) IsSuperAdmin() bool {
	return u.Role == "admin"
}

// IsUser checks if user has user role
func (u *User) IsUser() bool {
	return u.Role == "user"
}
