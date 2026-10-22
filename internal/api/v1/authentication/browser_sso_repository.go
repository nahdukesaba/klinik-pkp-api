package authentication

import (
	"context"
	"fmt"

	"klinik-pkp-api/internal/api/v1/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BrowserSSORepository struct {
	db *gorm.DB
}

func NewBrowserSSORepository(db *gorm.DB) *BrowserSSORepository {
	return &BrowserSSORepository{db: db}
}

func (r *BrowserSSORepository) CreateFlow(ctx context.Context, flow *BrowserSSOFlow) error {
	result := r.db.WithContext(ctx).Exec(`
INSERT INTO public.klinikpkp_sso_flow (state_hash, expires_at)
VALUES (?, ?)
`, flow.StateHash, flow.ExpiresAt)
	if result.Error != nil {
		return fmt.Errorf("create Klinik browser SSO flow: %w", result.Error)
	}
	return nil
}

func (r *BrowserSSORepository) ConsumeFlow(
	ctx context.Context,
	stateHash []byte,
) (*BrowserSSOFlow, error) {
	var flow BrowserSSOFlow
	result := r.db.WithContext(ctx).Raw(`
UPDATE public.klinikpkp_sso_flow
SET used_at = now()
WHERE state_hash = ?
  AND used_at IS NULL
  AND expires_at > now()
RETURNING state_hash, created_at, expires_at, used_at
`, stateHash).Scan(&flow)
	if result.Error != nil {
		return nil, fmt.Errorf("consume Klinik browser SSO flow: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &flow, nil
}

func (r *BrowserSSORepository) GetActiveByAppUserID(
	ctx context.Context,
	appUserID uuid.UUID,
) (*user.User, error) {
	var localUser user.User
	err := r.db.WithContext(ctx).
		Where("app_user_id = ? AND is_active = ?", appUserID, true).
		First(&localUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("resolve linked Klinik user: %w", err)
	}
	return &localUser, nil
}
