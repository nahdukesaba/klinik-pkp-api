package authentication

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/identity"
	"klinik-pkp-api/utils"

	"github.com/google/uuid"
)

type BrowserSSOFlow struct {
	StateHash []byte
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt    *time.Time
}

type BrowserSSOFlowStore interface {
	CreateFlow(context.Context, *BrowserSSOFlow) error
	ConsumeFlow(context.Context, []byte) (*BrowserSSOFlow, error)
}

type BrowserSSOCentralClient interface {
	ExchangeBrowserSSO(context.Context, string, string, string, string) (*utils.SSOSession, error)
	ResolveKlinikContext(context.Context, string) (*utils.SSOContext, error)
}

type BrowserSSOUserStore interface {
	GetByAppUserID(context.Context, uuid.UUID) (*identity.KlinikIdentity, error)
}

type BrowserSSOError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *BrowserSSOError) Error() string { return e.Message }

type BrowserSSOUser struct {
	ID        uuid.UUID `json:"id"`
	AppUserID uuid.UUID `json:"appUserId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	NIP       string    `json:"nip,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
}

type BrowserSSOCompletion struct {
	Status  string
	Session *utils.SSOSession
	User    *BrowserSSOUser
}

type BrowserSSOService struct {
	config  config.Config
	flows   BrowserSSOFlowStore
	central BrowserSSOCentralClient
	users   BrowserSSOUserStore
	now     func() time.Time
}

func NewBrowserSSOService(
	cfg config.Config,
	flows BrowserSSOFlowStore,
	central BrowserSSOCentralClient,
	users BrowserSSOUserStore,
) *BrowserSSOService {
	return &BrowserSSOService{config: cfg, flows: flows, central: central, users: users, now: time.Now}
}

func (s *BrowserSSOService) Start(ctx context.Context) (string, error) {
	key, err := s.key()
	if err != nil || s.flows == nil {
		return "", browserSSOUnavailable()
	}
	state, err := randomBrowserSSOValue(32)
	if err != nil {
		return "", browserSSOUnavailable()
	}
	if err := s.flows.CreateFlow(ctx, &BrowserSSOFlow{
		StateHash: browserSSOHash(state),
		ExpiresAt: s.now().Add(s.config.SSOFlowTTL),
	}); err != nil {
		return "", browserSSOUnavailable()
	}
	verifier := deriveKlinikPKCEVerifier(key, state)
	challengeHash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(challengeHash[:])

	authorizeURL, err := url.Parse(strings.TrimRight(s.config.SSOBaseURL, "/") + "/authorize")
	if err != nil {
		return "", browserSSOUnavailable()
	}
	query := authorizeURL.Query()
	query.Set("client_id", s.config.SSOClientID)
	query.Set("redirect_uri", s.config.SSORedirectURI)
	query.Set("response_type", "code")
	query.Set("state", state)
	query.Set("code_challenge", challenge)
	query.Set("code_challenge_method", "S256")
	authorizeURL.RawQuery = query.Encode()
	return authorizeURL.String(), nil
}

func (s *BrowserSSOService) LogoutURL(cancelReturn ...string) (string, error) {
	if !s.config.SSOBrowserEnabled {
		return "", browserSSOUnavailable()
	}
	if err := s.config.ValidateBrowserSSO(); err != nil {
		return "", browserSSOUnavailable()
	}
	logoutURL, err := url.Parse(strings.TrimRight(s.config.SSOBaseURL, "/") + "/logout")
	if err != nil {
		return "", browserSSOUnavailable()
	}
	query := logoutURL.Query()
	query.Set("client_id", s.config.SSOClientID)
	query.Set("post_logout_redirect_uri", s.config.SSOPostLogoutRedirectURI)
	if len(cancelReturn) > 0 && cancelReturn[0] != "" {
		if sanitized := sanitizeBrowserSSOCancelReturn(cancelReturn[0]); sanitized != "" {
			query.Set("cancel_return", sanitized)
		}
	}
	logoutURL.RawQuery = query.Encode()
	return logoutURL.String(), nil
}

func sanitizeBrowserSSOCancelReturn(value string) string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.HasPrefix(value, "/\\") ||
		strings.Contains(value, "\\") || strings.Contains(value, "#") {
		return ""
	}
	for i := 0; i < len(value); i++ {
		if value[i] <= ' ' || value[i] == 0x7f {
			return ""
		}
	}
	return value
}

func (s *BrowserSSOService) Complete(
	ctx context.Context,
	state string,
	code string,
	status string,
) (*BrowserSSOCompletion, error) {
	state = strings.TrimSpace(state)
	code = strings.TrimSpace(code)
	status = strings.ToLower(strings.TrimSpace(status))
	key, err := s.key()
	if err != nil || s.flows == nil || state == "" {
		return nil, invalidBrowserSSOState()
	}
	flow, err := s.flows.ConsumeFlow(ctx, browserSSOHash(state))
	if err != nil {
		return nil, browserSSOUnavailable()
	}
	if flow == nil {
		return nil, invalidBrowserSSOState()
	}

	switch status {
	case "pending":
		return &BrowserSSOCompletion{Status: "pending"}, nil
	case "denied", "rejected", "disabled":
		return nil, &BrowserSSOError{
			StatusCode: http.StatusForbidden,
			Code:       "KLINIK_ACCESS_DENIED",
			Message:    "Klinik PKP access is not available",
		}
	case "":
	default:
		return nil, &BrowserSSOError{
			StatusCode: http.StatusBadRequest,
			Code:       "INVALID_SSO_RESULT",
			Message:    "invalid SSO result",
		}
	}
	if code == "" {
		return nil, &BrowserSSOError{
			StatusCode: http.StatusBadRequest,
			Code:       "INVALID_SSO_RESULT",
			Message:    "authorization code is required",
		}
	}
	if s.central == nil || s.users == nil {
		return nil, browserSSOUnavailable()
	}
	session, err := s.central.ExchangeBrowserSSO(
		ctx,
		code,
		s.config.SSOClientID,
		s.config.SSORedirectURI,
		deriveKlinikPKCEVerifier(key, state),
	)
	if err != nil {
		if errors.Is(err, utils.ErrSSOUnauthorized) {
			return nil, &BrowserSSOError{
				StatusCode: http.StatusUnauthorized,
				Code:       "INVALID_AUTHORIZATION_CODE",
				Message:    "authorization code is invalid or expired",
			}
		}
		return nil, browserSSOUnavailable()
	}
	if session == nil || session.User.ID == uuid.Nil || strings.TrimSpace(session.AccessToken) == "" {
		return nil, browserSSOUnavailable()
	}
	contextResult, err := s.central.ResolveKlinikContext(ctx, session.AccessToken)
	if errors.Is(err, utils.ErrSSOUnavailable) {
		return nil, browserSSOUnavailable()
	}
	if err != nil || contextResult == nil || contextResult.AppUserID == uuid.Nil ||
		contextResult.AppUserID != session.User.ID || contextResult.Status != "approved" ||
		(contextResult.Role != "admin" && contextResult.Role != "user") {
		return nil, &BrowserSSOError{
			StatusCode: http.StatusForbidden,
			Code:       "KLINIK_ACCESS_DENIED",
			Message:    "Klinik PKP access is not available",
		}
	}
	canonical, err := s.users.GetByAppUserID(ctx, session.User.ID)
	if err != nil {
		return nil, browserSSOUnavailable()
	}
	if canonical == nil || !canonical.IsApproved() || canonical.AppUserID != session.User.ID ||
		canonical.Role != contextResult.Role {
		return nil, &BrowserSSOError{
			StatusCode: http.StatusForbidden,
			Code:       "KLINIK_ACCESS_DENIED",
			Message:    "Klinik PKP access is not available",
		}
	}
	phone := ""
	if canonical.Phone != nil {
		phone = *canonical.Phone
	}
	nip := ""
	if canonical.NIP != nil {
		nip = *canonical.NIP
	}
	return &BrowserSSOCompletion{
		Status:  "authenticated",
		Session: session,
		User: &BrowserSSOUser{
			ID: canonical.AppUserID, AppUserID: canonical.AppUserID,
			Name: canonical.FullName, Email: canonical.Email,
			Phone: phone, NIP: nip, Role: canonical.Role,
			Status: canonical.MembershipStatus,
		},
	}, nil
}

func (s *BrowserSSOService) key() ([]byte, error) {
	if err := s.config.ValidateBrowserSSO(); err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(s.config.SSOFlowHMACKey)
	if err != nil {
		key, err = base64.RawStdEncoding.DecodeString(s.config.SSOFlowHMACKey)
	}
	if err != nil || len(key) < 32 {
		return nil, errors.New("browser SSO is not configured")
	}
	return key, nil
}

func deriveKlinikPKCEVerifier(key []byte, state string) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte("klinik-pkp-pkce:" + state))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func randomBrowserSSOValue(size int) (string, error) {
	value := make([]byte, size)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func browserSSOHash(value string) []byte {
	hash := sha256.Sum256([]byte(value))
	return hash[:]
}

func browserSSOUnavailable() *BrowserSSOError {
	return &BrowserSSOError{
		StatusCode: http.StatusServiceUnavailable,
		Code:       "SSO_UNAVAILABLE",
		Message:    "browser SSO is unavailable",
	}
}

func invalidBrowserSSOState() *BrowserSSOError {
	return &BrowserSSOError{
		StatusCode: http.StatusUnauthorized,
		Code:       "INVALID_SSO_STATE",
		Message:    "invalid or expired SSO state",
	}
}
