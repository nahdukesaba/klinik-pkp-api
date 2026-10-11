package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSSOUnauthorized = errors.New("SSO token unauthorized")
	ErrSSOUnavailable  = errors.New("SSO service unavailable")
)

type SSOContext struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type SSOClient struct {
	baseURL string
	client  *http.Client
}

func NewSSOClient() *SSOClient {
	return &SSOClient{
		baseURL: strings.TrimRight(
			strings.TrimSpace(os.Getenv("SSO_BASE_URL")),
			"/",
		),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *SSOClient) ResolveKlinikContext(
	ctx context.Context,
	accessToken string,
) (*SSOContext, error) {
	accessToken = strings.TrimSpace(accessToken)

	if c.baseURL == "" {
		return nil, ErrSSOUnavailable
	}

	if accessToken == "" {
		return nil, ErrSSOUnauthorized
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/api/v1/auth/klinik-context",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create SSO context request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %v",
			ErrSSOUnavailable,
			err,
		)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Continue.

	case http.StatusUnauthorized,
		http.StatusForbidden:

		return nil, ErrSSOUnauthorized

	default:
		return nil, fmt.Errorf(
			"%w: SSO returned status %d",
			ErrSSOUnavailable,
			resp.StatusCode,
		)
	}

	var result SSOContext

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"%w: decode Klinik context: %v",
			ErrSSOUnavailable,
			err,
		)
	}

	if result.ID == uuid.Nil ||
		strings.TrimSpace(result.Email) == "" ||
		strings.TrimSpace(result.Role) == "" {

		return nil, ErrSSOUnauthorized
	}

	return &result, nil
}
