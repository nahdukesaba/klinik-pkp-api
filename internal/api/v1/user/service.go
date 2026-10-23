package user

import (
	"context"
	"errors"
	"strings"

	"klinik-pkp-api/internal/identity"
	"klinik-pkp-api/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrCanonicalIdentityUnavailable = errors.New("canonical Klinik identity unavailable")
	ErrDirectoryUserNotFound        = errors.New("canonical Klinik user not found")
	ErrInvalidDirectoryMutation     = errors.New("invalid canonical user mutation")
	ErrSelfAccessMutation           = errors.New("cannot remove current Klinik admin access")
	ErrLastAdminMutation            = errors.New("cannot remove the last approved Klinik admin")
	ErrProvisioningUnavailable      = errors.New("Central SSO provisioning unavailable")
)

type DirectoryQuery struct {
	Pagination utils.PaginationQuery
	Keyword    string
	Role       MembershipRole
	Status     MembershipStatus
	SortBy     string
	SortOrder  string
}

type DirectoryUpdate struct {
	FullName      string           `json:"fullName" validate:"required,ne=null,ne=NULL,ne=Null"`
	Phone         string           `json:"phone" validate:"-"`
	NIP           string           `json:"nip" validate:"-"`
	NIK           string           `json:"nik" validate:"-"`
	Alamat        string           `json:"alamat" validate:"-"`
	GelarDepan    string           `json:"gelarDepan" validate:"-"`
	GelarBelakang string           `json:"gelarBelakang" validate:"-"`
	Role          MembershipRole   `json:"role" validate:"required,oneof=admin user"`
	Status        MembershipStatus `json:"status" validate:"required,oneof=pending approved rejected disabled"`
}

type DirectoryCreate struct {
	Email string `json:"email" validate:"required,email"`
	DirectoryUpdate
}

type DirectoryStore interface {
	List(context.Context, DirectoryQuery) ([]DirectoryUser, int64, error)
	GetByAppUserID(context.Context, uuid.UUID) (*DirectoryUser, error)
	Update(context.Context, uuid.UUID, uuid.UUID, DirectoryUpdate) (*DirectoryUser, error)
	Disable(context.Context, uuid.UUID, uuid.UUID) error
}

type Provisioner interface {
	Provision(context.Context, string, DirectoryCreate) (uuid.UUID, error)
}

type Service struct {
	db          *gorm.DB
	identities  identity.Store
	directory   DirectoryStore
	provisioner Provisioner
	validator   *utils.Validator
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

func NewService(db *gorm.DB, centralBaseURL ...string) *Service {
	baseURL := ""
	if len(centralBaseURL) > 0 {
		baseURL = centralBaseURL[0]
	}
	return NewServiceWithStores(db, identity.NewRepository(db), NewCanonicalRepository(db), NewCentralProvisioner(baseURL))
}

func NewServiceWithIdentityStore(db *gorm.DB, identities identity.Store) *Service {
	return NewServiceWithStores(db, identities, NewCanonicalRepository(db), NewCentralProvisioner(""))
}

func NewServiceWithStores(db *gorm.DB, identities identity.Store, directory DirectoryStore, provisioner Provisioner) *Service {
	return &Service{db: db, identities: identities, directory: directory, provisioner: provisioner, validator: utils.NewValidator()}
}

func (s *Service) GetCurrentUser(ctx context.Context, appUserID uuid.UUID) (*CurrentUser, error) {
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
	phone, nip := "", ""
	if canonical.Phone != nil {
		phone = *canonical.Phone
	}
	if canonical.NIP != nil {
		nip = *canonical.NIP
	}
	return &CurrentUser{ID: canonical.AppUserID, AppUserID: canonical.AppUserID, AuthUserID: *canonical.AuthUserID, Name: canonical.FullName, Email: canonical.Email, Phone: phone, NIP: nip, Role: canonical.Role, Status: canonical.MembershipStatus, IsActive: true}, nil
}

func (s *Service) List(ctx context.Context, query DirectoryQuery) ([]DirectoryUser, int64, error) {
	if s == nil || s.directory == nil {
		return nil, 0, ErrCanonicalIdentityUnavailable
	}
	query.Keyword = strings.TrimSpace(query.Keyword)
	if query.Role != "" && !validRole(query.Role) {
		return nil, 0, ErrInvalidDirectoryMutation
	}
	if query.Status != "" && !validStatus(query.Status) {
		return nil, 0, ErrInvalidDirectoryMutation
	}
	return s.directory.List(ctx, query)
}

func (s *Service) Get(ctx context.Context, appUserID uuid.UUID) (*DirectoryUser, error) {
	if appUserID == uuid.Nil || s == nil || s.directory == nil {
		return nil, ErrDirectoryUserNotFound
	}
	result, err := s.directory.GetByAppUserID(ctx, appUserID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrDirectoryUserNotFound
	}
	return result, nil
}

func (s *Service) Create(ctx context.Context, authorization string, input DirectoryCreate) (*DirectoryUser, error) {
	normalizeCreate(&input)
	if input.Email == "" || input.FullName == "" || !validRole(input.Role) || !validStatus(input.Status) {
		return nil, ErrInvalidDirectoryMutation
	}
	if s == nil || s.provisioner == nil || s.directory == nil {
		return nil, ErrProvisioningUnavailable
	}
	appUserID, err := s.provisioner.Provision(ctx, authorization, input)
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, appUserID)
}

func (s *Service) Update(ctx context.Context, actorAppUserID, targetAppUserID uuid.UUID, input DirectoryUpdate) (*DirectoryUser, error) {
	normalizeUpdate(&input)
	if actorAppUserID == uuid.Nil || targetAppUserID == uuid.Nil || input.FullName == "" || !validRole(input.Role) || !validStatus(input.Status) {
		return nil, ErrInvalidDirectoryMutation
	}
	return s.directory.Update(ctx, actorAppUserID, targetAppUserID, input)
}

func (s *Service) Disable(ctx context.Context, actorAppUserID, targetAppUserID uuid.UUID) error {
	if actorAppUserID == uuid.Nil || targetAppUserID == uuid.Nil {
		return ErrInvalidDirectoryMutation
	}
	return s.directory.Disable(ctx, actorAppUserID, targetAppUserID)
}

func normalizeCreate(input *DirectoryCreate) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	normalizeUpdate(&input.DirectoryUpdate)
}
func normalizeUpdate(input *DirectoryUpdate) {
	input.FullName = strings.TrimSpace(input.FullName)
	input.Phone = strings.TrimSpace(input.Phone)
	input.NIP = strings.TrimSpace(input.NIP)
	input.NIK = strings.TrimSpace(input.NIK)
	input.Alamat = strings.TrimSpace(input.Alamat)
	input.GelarDepan = strings.TrimSpace(input.GelarDepan)
	input.GelarBelakang = strings.TrimSpace(input.GelarBelakang)
	input.Role = MembershipRole(strings.ToLower(strings.TrimSpace(string(input.Role))))
	input.Status = MembershipStatus(strings.ToLower(strings.TrimSpace(string(input.Status))))
}
func validRole(role MembershipRole) bool { return role == RoleAdmin || role == RoleUser }
func validStatus(status MembershipStatus) bool {
	return status == StatusPending || status == StatusApproved || status == StatusRejected || status == StatusDisabled
}
