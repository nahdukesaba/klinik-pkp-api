package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"size:255;not null" json:"name"`
	Email    string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string    `gorm:"size:255;not null" json:"-"`
	Phone    string    `gorm:"size:20;uniqueIndex" json:"phone,omitempty"`
	NIP      string    `gorm:"column:nip;size:50;uniqueIndex:idx_user_nip" json:"nip,omitempty"`
	Role     string    `gorm:"type:varchar(50);not null" json:"role"`
	IsActive bool      `gorm:"default:true" json:"is_active"`

	CreatedBy uuid.UUID `gorm:"index;" json:"created_by,omitempty"`
	UpdatedBy uuid.UUID `gorm:"index;" json:"updated_by,omitempty"`
	DeletedBy uuid.UUID `gorm:"index;" json:"deleted_by,omitempty"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "user"
}

// HASH PASSWORD BEFORE CREATING USER
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
