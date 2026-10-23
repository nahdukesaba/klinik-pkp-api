package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID  `gorm:"primaryKey" json:"id"`
	AppUserID *uuid.UUID `gorm:"column:app_user_id;type:uuid" json:"-"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	Email     string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"size:255;not null" json:"-"`
	Phone     string     `gorm:"size:20;uniqueIndex" json:"phone,omitempty"`
	NIP       string     `gorm:"column:nip;size:50;uniqueIndex:idx_user_nip" json:"nip,omitempty"`
	Role      string     `gorm:"type:varchar(50);not null" json:"role"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`

	CreatedBy uuid.UUID `gorm:"index;" json:"created_by,omitempty"`
	UpdatedBy uuid.UUID `gorm:"index;" json:"updated_by,omitempty"`
	DeletedBy uuid.UUID `gorm:"index;" json:"deleted_by,omitempty"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type MembershipRole string
type MembershipStatus string

const (
	RoleAdmin      MembershipRole   = "admin"
	RoleUser       MembershipRole   = "user"
	StatusPending  MembershipStatus = "pending"
	StatusApproved MembershipStatus = "approved"
	StatusRejected MembershipStatus = "rejected"
	StatusDisabled MembershipStatus = "disabled"
)

// DirectoryUser is the canonical Klinik admin projection. It deliberately has
// no local Klinik user ID, password, local role, or is_active field.
type DirectoryUser struct {
	AppUserID     uuid.UUID        `gorm:"column:app_user_id" json:"appUserId"`
	AuthUserID    *uuid.UUID       `gorm:"column:auth_user_id" json:"authUserId,omitempty"`
	FullName      string           `gorm:"column:full_name" json:"fullName"`
	Email         string           `gorm:"column:email" json:"email"`
	Phone         *string          `gorm:"column:phone" json:"phone,omitempty"`
	NIP           *string          `gorm:"column:nip" json:"nip,omitempty"`
	NIK           *string          `gorm:"column:nik" json:"nik,omitempty"`
	Alamat        *string          `gorm:"column:alamat" json:"alamat,omitempty"`
	GelarDepan    *string          `gorm:"column:gelar_depan" json:"gelarDepan,omitempty"`
	GelarBelakang *string          `gorm:"column:gelar_belakang" json:"gelarBelakang,omitempty"`
	Role          MembershipRole   `gorm:"column:role" json:"role"`
	Status        MembershipStatus `gorm:"column:status" json:"status"`
	CreatedAt     time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

// func (User) TableName() string {
// 	return "user"
// }

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
