package utils

import (
	"bytes"
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

type SSOSessionUser struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"fullName"`
}

type SSOSession struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	TokenType    string         `json:"tokenType"`
	ExpiresIn    int            `json:"expiresIn"`
	ExpiresAt    int64          `json:"expiresAt"`
	User         SSOSessionUser `json:"user"`
}

type SSOClient struct {
	baseURL string
	client  *http.Client
}

func NewSSOClient() *SSOClient {
	return NewSSOClientWithBaseURL(
		strings.TrimRight(
			strings.TrimSpace(os.Getenv("SSO_BASE_URL")),
			"/",
		),
	)
}

func NewSSOClientWithBaseURL(baseURL string) *SSOClient {
	return &SSOClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *SSOClient) ExchangeBrowserSSO(
	ctx context.Context,
	code string,
	clientID string,
	redirectURI string,
	codeVerifier string,
) (*SSOSession, error) {
	if c.baseURL == "" {
		return nil, ErrSSOUnavailable
	}
	payload, err := json.Marshal(map[string]string{
		"code":         strings.TrimSpace(code),
		"clientId":     strings.TrimSpace(clientID),
		"redirectUri":  strings.TrimSpace(redirectURI),
		"codeVerifier": strings.TrimSpace(codeVerifier),
	})
	if err != nil {
		return nil, fmt.Errorf("create SSO exchange payload: %w", err)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/v1/auth/sso/exchange",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("create SSO exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSSOUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized ||
			resp.StatusCode == http.StatusForbidden {
			return nil, ErrSSOUnauthorized
		}
		return nil, fmt.Errorf("%w: SSO returned status %d", ErrSSOUnavailable, resp.StatusCode)
	}
	var session SSOSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("%w: decode SSO exchange: %v", ErrSSOUnavailable, err)
	}
	if strings.TrimSpace(session.AccessToken) == "" || strings.TrimSpace(session.RefreshToken) == "" ||
		session.User.ID == uuid.Nil || strings.TrimSpace(session.User.Email) == "" {
		return nil, fmt.Errorf("%w: invalid SSO exchange response", ErrSSOUnavailable)
	}
	return &session, nil
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
		!IsKlinikRole(strings.ToLower(strings.TrimSpace(result.Role))) {

		return nil, ErrSSOUnauthorized
	}
	result.Role = strings.ToLower(strings.TrimSpace(result.Role))

	return &result, nil
}
