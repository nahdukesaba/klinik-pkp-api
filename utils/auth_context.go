package utils

import (
	"context"
	"errors"
	"strings"

	"klinik-pkp-api/internal/identity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginSource string

const (
	LoginSourceDirect LoginSource = "direct"
	LoginSourceSSO    LoginSource = "sso"
)

var (
	ErrKlinikAuthorizationDenied      = errors.New("Klinik authorization denied")
	ErrKlinikAuthorizationUnavailable = errors.New("Klinik authorization unavailable")
)

// AuthContext is the trusted authorization context used by downstream handlers.
// Its role is always loaded from the Central app_memberships table.
type AuthContext struct {
	AppUserID        uuid.UUID   `json:"appUserId"`
	AuthUserID       *uuid.UUID  `json:"authUserId,omitempty"`
	Email            string      `json:"email"`
	FullName         string      `json:"fullName"`
	Phone            *string     `json:"phone,omitempty"`
	NIP              *string     `json:"nip,omitempty"`
	Role             string      `json:"role"`
	MembershipStatus string      `json:"membershipStatus"`
	LoginSource      LoginSource `json:"loginSource"`
	LocalUserID      *uuid.UUID  `json:"localUserId,omitempty"`
}

type KlinikAuthorizationResolver interface {
	ResolveKlinikAuthorization(context.Context, uuid.UUID) (*AuthContext, error)
	ResolveLegacyKlinikAuthorization(context.Context, uuid.UUID) (*AuthContext, error)
}

type databaseKlinikAuthorizationResolver struct {
	db         *gorm.DB
	identities identity.Store
}

func NewKlinikAuthorizationResolver(db *gorm.DB) KlinikAuthorizationResolver {
	return &databaseKlinikAuthorizationResolver{
		db:         db,
		identities: identity.NewRepository(db),
	}
}

func (r *databaseKlinikAuthorizationResolver) ResolveKlinikAuthorization(
	ctx context.Context,
	appUserID uuid.UUID,
) (*AuthContext, error) {
	if r == nil || r.identities == nil || appUserID == uuid.Nil {
		return nil, ErrKlinikAuthorizationDenied
	}
	canonical, err := r.identities.GetByAppUserID(ctx, appUserID)
	if err != nil {
		return nil, errors.Join(ErrKlinikAuthorizationUnavailable, err)
	}
	trusted, err := authorizedKlinikContext(canonical)
	if err != nil {
		return nil, err
	}

	// Local ID is retained only as best-effort compatibility for legacy CRUD.
	// A missing local row must never deny canonical authentication.
	trusted.LocalUserID = r.findOptionalLocalUserID(ctx, appUserID)
	return trusted, nil
}

func (r *databaseKlinikAuthorizationResolver) ResolveLegacyKlinikAuthorization(
	ctx context.Context,
	localUserID uuid.UUID,
) (*AuthContext, error) {
	if r == nil || r.db == nil || localUserID == uuid.Nil {
		return nil, ErrKlinikAuthorizationDenied
	}

	var result struct {
		AppUserID uuid.UUID `gorm:"column:app_user_id"`
	}

	query := r.db.WithContext(ctx).Raw(`
SELECT
    ku.app_user_id
FROM public.klinikpkp_user AS ku
WHERE ku.id = ?
  AND ku.app_user_id IS NOT NULL
  AND ku.is_active = TRUE
  AND ku.deleted_at IS NULL
LIMIT 1`, localUserID).Scan(&result)
	if query.Error != nil {
		return nil, errors.Join(ErrKlinikAuthorizationUnavailable, query.Error)
	}
	if query.RowsAffected != 1 || result.AppUserID == uuid.Nil {
		return nil, ErrKlinikAuthorizationDenied
	}

	trusted, err := r.ResolveKlinikAuthorization(ctx, result.AppUserID)
	if err != nil {
		return nil, err
	}
	trusted.LocalUserID = &localUserID
	return trusted, nil
}

func (r *databaseKlinikAuthorizationResolver) findOptionalLocalUserID(
	ctx context.Context,
	appUserID uuid.UUID,
) *uuid.UUID {
	if r == nil || r.db == nil {
		return nil
	}
	var result struct {
		ID uuid.UUID `gorm:"column:id"`
	}
	query := r.db.WithContext(ctx).Raw(`
SELECT id
FROM public.klinikpkp_user
WHERE app_user_id = ?
  AND deleted_at IS NULL
LIMIT 1`, appUserID).Scan(&result)
	if query.Error != nil || query.RowsAffected != 1 || result.ID == uuid.Nil {
		return nil
	}
	return &result.ID
}

func authorizedKlinikContext(canonical *identity.KlinikIdentity) (*AuthContext, error) {
	if canonical == nil || !canonical.IsApproved() {
		return nil, ErrKlinikAuthorizationDenied
	}
	return &AuthContext{
		AppUserID:        canonical.AppUserID,
		AuthUserID:       canonical.AuthUserID,
		Email:            strings.TrimSpace(canonical.Email),
		FullName:         strings.TrimSpace(canonical.FullName),
		Phone:            canonical.Phone,
		NIP:              canonical.NIP,
		Role:             canonical.Role,
		MembershipStatus: canonical.MembershipStatus,
	}, nil
}

func IsKlinikRole(role string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	return role == "admin" || role == "user"
}

func TrustedAuthContext(ctx *fiber.Ctx) (*AuthContext, bool) {
	trusted, ok := ctx.Locals("auth_context").(*AuthContext)
	return trusted, ok && trusted != nil
}
