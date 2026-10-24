package authentication

import (
	"context"
	"encoding/hex"
	"fmt"

	"gorm.io/gorm"
)

type BrowserSSORepository struct {
	db *gorm.DB
}

func NewBrowserSSORepository(db *gorm.DB) *BrowserSSORepository {
	return &BrowserSSORepository{db: db}
}

func (r *BrowserSSORepository) CreateFlow(
	ctx context.Context,
	flow *BrowserSSOFlow,
) error {
	stateHashHex := hex.EncodeToString(flow.StateHash)

	result := r.db.WithContext(ctx).Exec(`
INSERT INTO public.klinikpkp_sso_flow (state_hash, expires_at)
VALUES (decode(?, 'hex'), ?)
`, stateHashHex, flow.ExpiresAt)

	if result.Error != nil {
		return fmt.Errorf(
			"create Klinik browser SSO flow: %w",
			result.Error,
		)
	}

	return nil
}

func (r *BrowserSSORepository) ConsumeFlow(
	ctx context.Context,
	stateHash []byte,
) (*BrowserSSOFlow, error) {
	var flow BrowserSSOFlow

	stateHashHex := hex.EncodeToString(stateHash)

	result := r.db.WithContext(ctx).Raw(`
UPDATE public.klinikpkp_sso_flow
SET used_at = now()
WHERE state_hash = decode(?, 'hex')
  AND used_at IS NULL
  AND expires_at > now()
RETURNING state_hash, created_at, expires_at, used_at
`, stateHashHex).Scan(&flow)

	if result.Error != nil {
		return nil, fmt.Errorf(
			"consume Klinik browser SSO flow: %w",
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &flow, nil
}
