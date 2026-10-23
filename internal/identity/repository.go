package identity

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	GetByAppUserID(context.Context, uuid.UUID) (*KlinikIdentity, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByAppUserID(
	ctx context.Context,
	appUserID uuid.UUID,
) (*KlinikIdentity, error) {
	if r == nil || r.db == nil || appUserID == uuid.Nil {
		return nil, nil
	}

	var result KlinikIdentity
	query := r.db.WithContext(ctx).Raw(`
SELECT
    app_user.id AS app_user_id,
    app_user.auth_user_id,
    app_user.email,
    app_user.full_name,
    app_user.phone,
    pegawai.nip,
    membership.role,
    membership.status AS membership_status,
    application.code AS application_code
FROM public.app_users AS app_user
JOIN public.app_memberships AS membership
    ON membership.app_user_id = app_user.id
   AND membership.deleted_at IS NULL
JOIN public.applications AS application
    ON application.id = membership.application_id
   AND application.code = ?
   AND application.is_active = TRUE
LEFT JOIN public.app_pegawai AS pegawai
    ON pegawai.app_user_id = app_user.id
   AND pegawai.deleted_at IS NULL
WHERE app_user.id = ?
  AND app_user.auth_user_id IS NOT NULL
  AND app_user.deleted_at IS NULL
LIMIT 1`, KlinikApplicationCode, appUserID).Scan(&result)
	if query.Error != nil {
		return nil, fmt.Errorf("get canonical Klinik identity: %w", query.Error)
	}
	if query.RowsAffected != 1 || result.AppUserID == uuid.Nil {
		return nil, nil
	}
	result.Normalize()
	return &result, nil
}
