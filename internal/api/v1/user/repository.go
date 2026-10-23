package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const canonicalDirectoryFrom = `
FROM public.app_users AS au
JOIN public.app_memberships AS am
  ON am.app_user_id = au.id AND am.deleted_at IS NULL
JOIN public.applications AS a
  ON a.id = am.application_id AND a.code = ? AND a.is_active = TRUE
LEFT JOIN public.app_pegawai AS ap
  ON ap.app_user_id = au.id AND ap.deleted_at IS NULL
WHERE au.deleted_at IS NULL`

const canonicalDirectoryColumns = `
SELECT au.id AS app_user_id, au.auth_user_id, au.full_name, au.email, au.phone,
       ap.nip, ap.nik, ap.alamat, ap.gelar_depan, ap.gelar_belakang,
       am.role, am.status, am.created_at, am.updated_at`

const identityApplicationCode = "klinik_pkp"

type CanonicalRepository struct{ db *gorm.DB }

func NewCanonicalRepository(db *gorm.DB) *CanonicalRepository { return &CanonicalRepository{db: db} }

func (r *CanonicalRepository) List(ctx context.Context, query DirectoryQuery) ([]DirectoryUser, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, ErrCanonicalIdentityUnavailable
	}
	where, args := canonicalDirectoryFrom, []any{identityApplicationCode}
	if query.Keyword != "" {
		where += ` AND (au.full_name ILIKE ? OR au.email ILIKE ? OR COALESCE(au.phone, '') ILIKE ? OR COALESCE(ap.nip, '') ILIKE ?)`
		term := "%" + query.Keyword + "%"
		args = append(args, term, term, term, term)
	}
	if query.Role != "" {
		where += " AND am.role = ?"
		args = append(args, query.Role)
	}
	if query.Status != "" {
		where += " AND am.status = ?"
		args = append(args, query.Status)
	}
	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(*) "+where, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count canonical Klinik directory: %w", err)
	}
	sortColumn := map[string]string{"full_name": "au.full_name", "email": "au.email", "role": "am.role", "status": "am.status", "created_at": "am.created_at", "updated_at": "am.updated_at"}[query.SortBy]
	if sortColumn == "" {
		sortColumn = "au.full_name"
	}
	sortOrder := "ASC"
	if strings.EqualFold(query.SortOrder, "desc") {
		sortOrder = "DESC"
	}
	dataArgs := append(append([]any{}, args...), query.Pagination.Limit, query.Pagination.Offset())
	var users []DirectoryUser
	statement := canonicalDirectoryColumns + where + " ORDER BY " + sortColumn + " " + sortOrder + ", au.id ASC LIMIT ? OFFSET ?"
	if err := r.db.WithContext(ctx).Raw(statement, dataArgs...).Scan(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list canonical Klinik directory: %w", err)
	}
	return users, total, nil
}

func (r *CanonicalRepository) GetByAppUserID(ctx context.Context, appUserID uuid.UUID) (*DirectoryUser, error) {
	if r == nil || r.db == nil || appUserID == uuid.Nil {
		return nil, nil
	}
	var result DirectoryUser
	query := r.db.WithContext(ctx).Raw(canonicalDirectoryColumns+canonicalDirectoryFrom+" AND au.id = ? LIMIT 1", identityApplicationCode, appUserID).Scan(&result)
	if query.Error != nil {
		return nil, fmt.Errorf("get canonical Klinik directory user: %w", query.Error)
	}
	if query.RowsAffected != 1 || result.AppUserID == uuid.Nil {
		return nil, nil
	}
	return &result, nil
}

func (r *CanonicalRepository) Update(ctx context.Context, actorAppUserID, targetAppUserID uuid.UUID, input DirectoryUpdate) (*DirectoryUser, error) {
	if err := r.mutate(ctx, actorAppUserID, targetAppUserID, input, false); err != nil {
		return nil, err
	}
	return r.GetByAppUserID(ctx, targetAppUserID)
}

func (r *CanonicalRepository) Disable(ctx context.Context, actorAppUserID, targetAppUserID uuid.UUID) error {
	return r.mutate(ctx, actorAppUserID, targetAppUserID, DirectoryUpdate{Status: StatusDisabled}, true)
}

func (r *CanonicalRepository) mutate(ctx context.Context, actorAppUserID, targetAppUserID uuid.UUID, input DirectoryUpdate, disableOnly bool) error {
	if r == nil || r.db == nil {
		return ErrCanonicalIdentityUnavailable
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(hashtext(?))", "klinik_pkp_admin_guard").Error; err != nil {
			return err
		}
		var current struct {
			Role   MembershipRole
			Status MembershipStatus
		}
		row := tx.Raw(`SELECT am.role, am.status FROM public.app_memberships am JOIN public.applications a ON a.id = am.application_id WHERE am.app_user_id = ? AND a.code = ? AND a.is_active = TRUE AND am.deleted_at IS NULL FOR UPDATE`, targetAppUserID, identityApplicationCode).Scan(&current)
		if row.Error != nil {
			return row.Error
		}
		if row.RowsAffected != 1 {
			return ErrDirectoryUserNotFound
		}
		newRole, newStatus := input.Role, input.Status
		if disableOnly {
			newRole, newStatus = current.Role, StatusDisabled
		}
		removesAdmin := current.Role == RoleAdmin && current.Status == StatusApproved && (newRole != RoleAdmin || newStatus != StatusApproved)
		if targetAppUserID == actorAppUserID && removesAdmin {
			return ErrSelfAccessMutation
		}
		if removesAdmin {
			var count int64
			if err := tx.Raw(`SELECT COUNT(*) FROM public.app_memberships am JOIN public.applications a ON a.id = am.application_id WHERE a.code = ? AND a.is_active = TRUE AND am.role = 'admin' AND am.status = 'approved' AND am.deleted_at IS NULL`, identityApplicationCode).Scan(&count).Error; err != nil {
				return err
			}
			if count <= 1 {
				return ErrLastAdminMutation
			}
		}
		if !disableOnly {
			result := tx.Exec(`UPDATE public.app_users SET full_name = ?, phone = NULLIF(?, ''), updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`, input.FullName, input.Phone, targetAppUserID)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return ErrDirectoryUserNotFound
			}
			if err := tx.Exec(`INSERT INTO public.app_pegawai (app_user_id, nip, nik, alamat, gelar_depan, gelar_belakang) VALUES (?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, '')) ON CONFLICT (app_user_id) DO UPDATE SET nip=EXCLUDED.nip, nik=EXCLUDED.nik, alamat=EXCLUDED.alamat, gelar_depan=EXCLUDED.gelar_depan, gelar_belakang=EXCLUDED.gelar_belakang, deleted_at=NULL, updated_at=NOW()`, targetAppUserID, input.NIP, input.NIK, input.Alamat, input.GelarDepan, input.GelarBelakang).Error; err != nil {
				return err
			}
		}
		result := tx.Exec(`UPDATE public.app_memberships AS am SET role = ?, status = ?, updated_at = NOW() FROM public.applications AS a WHERE am.application_id = a.id AND am.app_user_id = ? AND a.code = ? AND a.is_active = TRUE AND am.deleted_at IS NULL`, newRole, newStatus, targetAppUserID, identityApplicationCode)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrDirectoryUserNotFound
		}
		return nil
	})
}
