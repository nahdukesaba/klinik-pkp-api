package user

import (
	"context"
	"errors"

	"klinik-pkp-api/internal/identity"
	"klinik-pkp-api/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db         *gorm.DB
	identities identity.Store
	validator  *utils.Validator
}

type CurrentUser struct {
	ID         uuid.UUID `json:"id"`
	AppUserID  uuid.UUID `json:"appUserId"`
	AuthUserID uuid.UUID `json:"authUserId"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone,omitempty"`
	NIP        string    `json:"nip,omitempty"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	IsActive   bool      `json:"is_active"`
}

var ErrCanonicalIdentityUnavailable = errors.New("canonical Klinik identity unavailable")

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
	return NewServiceWithIdentityStore(db, identity.NewRepository(db))
}

func NewServiceWithIdentityStore(db *gorm.DB, identities identity.Store) *Service {
	return &Service{
		db:         db,
		identities: identities,
		validator:  utils.NewValidator(),
	}
}

func (s *Service) GetCurrentUser(
	ctx context.Context,
	appUserID uuid.UUID,
) (*CurrentUser, error) {
	if s == nil || s.identities == nil || appUserID == uuid.Nil {
		return nil, ErrCanonicalIdentityUnavailable
	}
	canonical, err := s.identities.GetByAppUserID(ctx, appUserID)
	if err != nil {
		return nil, err
	}
	if canonical == nil || canonical.AppUserID != appUserID || !canonical.IsApproved() {
		return nil, ErrCanonicalIdentityUnavailable
	}
	phone := ""
	if canonical.Phone != nil {
		phone = *canonical.Phone
	}
	nip := ""
	if canonical.NIP != nil {
		nip = *canonical.NIP
	}
	return &CurrentUser{
		ID:         canonical.AppUserID,
		AppUserID:  canonical.AppUserID,
		AuthUserID: *canonical.AuthUserID,
		Name:       canonical.FullName,
		Email:      canonical.Email,
		Phone:      phone,
		NIP:        nip,
		Role:       canonical.Role,
		Status:     canonical.MembershipStatus,
		IsActive:   true,
	}, nil
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
