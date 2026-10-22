package utils

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const klinikApplicationCode = "klinik_pkp"

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
	LocalUserID uuid.UUID   `json:"localUserId"`
	AppUserID   uuid.UUID   `json:"appUserId"`
	Email       string      `json:"email"`
	Role        string      `json:"role"`
	LoginSource LoginSource `json:"loginSource"`
}

type KlinikAuthorizationResolver interface {
	ResolveKlinikAuthorization(context.Context, uuid.UUID) (*AuthContext, error)
}

type databaseKlinikAuthorizationResolver struct {
	db *gorm.DB
}

func NewKlinikAuthorizationResolver(db *gorm.DB) KlinikAuthorizationResolver {
	return &databaseKlinikAuthorizationResolver{db: db}
}

func (r *databaseKlinikAuthorizationResolver) ResolveKlinikAuthorization(
	ctx context.Context,
	localUserID uuid.UUID,
) (*AuthContext, error) {
	if r.db == nil || localUserID == uuid.Nil {
		return nil, ErrKlinikAuthorizationDenied
	}

	var result struct {
		LocalUserID uuid.UUID `gorm:"column:local_user_id"`
		AppUserID   uuid.UUID `gorm:"column:app_user_id"`
		Email       string    `gorm:"column:email"`
		Role        string    `gorm:"column:role"`
		Status      string    `gorm:"column:status"`
	}

	query := r.db.WithContext(ctx).Raw(`
SELECT
    ku.id AS local_user_id,
    ku.app_user_id,
    ku.email,
    membership.role,
    membership.status
FROM public.klinikpkp_user AS ku
JOIN public.app_users AS app_user
    ON app_user.id = ku.app_user_id
   AND app_user.deleted_at IS NULL
JOIN public.applications AS application
    ON application.code = ?
   AND application.is_active = TRUE
JOIN public.app_memberships AS membership
    ON membership.app_user_id = app_user.id
   AND membership.application_id = application.id
   AND membership.deleted_at IS NULL
WHERE ku.id = ?
  AND ku.app_user_id IS NOT NULL
  AND ku.is_active = TRUE
  AND ku.deleted_at IS NULL
LIMIT 1`, klinikApplicationCode, localUserID).Scan(&result)
	if query.Error != nil {
		return nil, errors.Join(ErrKlinikAuthorizationUnavailable, query.Error)
	}
	if query.RowsAffected != 1 || result.LocalUserID == uuid.Nil || result.AppUserID == uuid.Nil {
		return nil, ErrKlinikAuthorizationDenied
	}

	return authorizedKlinikContext(
		result.LocalUserID,
		result.AppUserID,
		result.Email,
		result.Role,
		result.Status,
	)
}

func authorizedKlinikContext(
	localUserID uuid.UUID,
	appUserID uuid.UUID,
	email string,
	role string,
	status string,
) (*AuthContext, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if localUserID == uuid.Nil || appUserID == uuid.Nil || status != "approved" || !IsKlinikRole(role) {
		return nil, ErrKlinikAuthorizationDenied
	}
	return &AuthContext{
		LocalUserID: localUserID,
		AppUserID:   appUserID,
		Email:       strings.TrimSpace(email),
		Role:        role,
	}, nil
}

func IsKlinikRole(role string) bool {
	return role == "admin" || role == "user"
}

func TrustedAuthContext(ctx *fiber.Ctx) (*AuthContext, bool) {
	trusted, ok := ctx.Locals("auth_context").(*AuthContext)
	return trusted, ok && trusted != nil
}
