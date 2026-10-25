package authentication

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/identity"
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	testAppUserID = uuid.MustParse("11111111-1111-4111-8111-111111111111")
)

type fakeBrowserSSOFlowStore struct {
	flow     *BrowserSSOFlow
	consumed bool
}

func (f *fakeBrowserSSOFlowStore) CreateFlow(_ context.Context, flow *BrowserSSOFlow) error {
	copyFlow := *flow
	copyFlow.StateHash = append([]byte(nil), flow.StateHash...)
	f.flow = &copyFlow
	return nil
}

func (f *fakeBrowserSSOFlowStore) ConsumeFlow(_ context.Context, hash []byte) (*BrowserSSOFlow, error) {
	if f.consumed || f.flow == nil || string(f.flow.StateHash) != string(hash) ||
		!f.flow.ExpiresAt.After(time.Now()) {
		return nil, nil
	}
	f.consumed = true
	return f.flow, nil
}

type fakeBrowserSSOCentral struct {
	session       *utils.SSOSession
	context       *utils.SSOContext
	exchangeErr   error
	contextErr    error
	exchangeCalls int
	contextCalls  int
	code          string
	clientID      string
	redirectURI   string
	codeVerifier  string
	contextToken  string
}

func (f *fakeBrowserSSOCentral) ExchangeBrowserSSO(
	_ context.Context,
	code string,
	clientID string,
	redirectURI string,
	codeVerifier string,
) (*utils.SSOSession, error) {
	f.exchangeCalls++
	f.code = code
	f.clientID = clientID
	f.redirectURI = redirectURI
	f.codeVerifier = codeVerifier
	return f.session, f.exchangeErr
}

func (f *fakeBrowserSSOCentral) ResolveKlinikContext(
	_ context.Context,
	accessToken string,
) (*utils.SSOContext, error) {
	f.contextCalls++
	f.contextToken = accessToken
	return f.context, f.contextErr
}

type fakeBrowserSSOUsers struct {
	canonical *identity.KlinikIdentity
	err       error
	requested uuid.UUID
}

func (f *fakeBrowserSSOUsers) GetByAppUserID(
	_ context.Context,
	appUserID uuid.UUID,
) (*identity.KlinikIdentity, error) {
	f.requested = appUserID
	return f.canonical, f.err
}

func browserSSOTestConfig() config.Config {
	return config.Config{
		SSOBrowserEnabled:        true,
		SSOBaseURL:               "http://127.0.0.1:18082",
		SSOClientID:              "klinik-pkp-web",
		SSOIntrospectionSecret:   "klinik-test-introspection-secret-00001",
		SSORedirectURI:           "http://localhost:3000/auth/sso/callback",
		SSOPostLogoutRedirectURI: "http://localhost:3000/login",
		SSOFlowTTL:               5 * time.Minute,
		SSOFlowHMACKey: base64.StdEncoding.EncodeToString(
			[]byte("01234567890123456789012345678901"),
		),
	}
}

func browserSSOTestService() (*BrowserSSOService, *fakeBrowserSSOFlowStore, *fakeBrowserSSOCentral, *fakeBrowserSSOUsers) {
	flows := &fakeBrowserSSOFlowStore{}
	central := &fakeBrowserSSOCentral{
		session: &utils.SSOSession{
			SID:         "11111111-1111-4111-8111-111111111111",
			AccessToken: "central-access", RefreshToken: "central-refresh", TokenType: "bearer",
			ExpiresIn: 3600, ExpiresAt: 9999999999,
			User: utils.SSOSessionUser{ID: testAppUserID, Email: "person@example.com"},
		},
		context: &utils.SSOContext{
			ID: testAppUserID, AppUserID: testAppUserID,
			AuthUserID: uuid.MustParse("33333333-3333-4333-8333-333333333333"),
			Email:      "person@example.com", Role: "user", Status: "approved",
		},
	}
	authUserID := central.context.AuthUserID
	phone := "08123456789"
	nip := "198001012010011001"
	users := &fakeBrowserSSOUsers{canonical: &identity.KlinikIdentity{
		AppUserID: testAppUserID, AuthUserID: &authUserID,
		Email: "person@example.com", FullName: "Pegawai Klinik",
		Phone: &phone, NIP: &nip, Role: "user", MembershipStatus: "approved",
		ApplicationCode: identity.KlinikApplicationCode,
	}}
	return NewBrowserSSOService(browserSSOTestConfig(), flows, central, users), flows, central, users
}

func startBrowserSSO(t *testing.T, service *BrowserSSOService) (string, *url.URL) {
	t.Helper()
	authorizeURL, err := service.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(authorizeURL)
	if err != nil {
		t.Fatal(err)
	}
	return parsed.Query().Get("state"), parsed
}

func TestBrowserSSOStartStoresOnlyHashAndBuildsTrustedPKCERequest(t *testing.T) {
	service, flows, _, _ := browserSSOTestService()
	state, parsed := startBrowserSSO(t, service)
	query := parsed.Query()
	if state == "" || query.Get("code_challenge") == "" || query.Get("code_challenge_method") != "S256" {
		t.Fatalf("secure authorize parameters missing: %s", parsed.RawQuery)
	}
	if parsed.Scheme != "http" || parsed.Host != "127.0.0.1:18082" || parsed.Path != "/authorize" ||
		query.Get("client_id") != "klinik-pkp-web" ||
		query.Get("redirect_uri") != "http://localhost:3000/auth/sso/callback" ||
		query.Get("response_type") != "code" {
		t.Fatalf("unexpected trusted authorize request: %q", parsed.String())
	}
	if query.Get("code_verifier") != "" || query.Get("access_token") != "" ||
		query.Get("refresh_token") != "" || query.Get("application") != "" {
		t.Fatal("authorize URL leaked secret or accepted application context")
	}
	if flows.flow == nil || string(flows.flow.StateHash) == state ||
		string(flows.flow.StateHash) != string(browserSSOHash(state)) {
		t.Fatal("flow did not persist only the state hash")
	}
	key, err := service.key()
	if err != nil {
		t.Fatal(err)
	}
	verifier := deriveKlinikPKCEVerifier(key, state)
	challengeHash := browserSSOHash(verifier)
	if query.Get("code_challenge") != base64.RawURLEncoding.EncodeToString(challengeHash) {
		t.Fatal("authorize request did not use the derived S256 challenge")
	}

	secondState, _ := startBrowserSSO(t, service)
	if secondState == state {
		t.Fatal("browser SSO start reused state")
	}
}

func TestBrowserSSOFlowRejectsUnknownExpiredAndReplay(t *testing.T) {
	service, flows, central, _ := browserSSOTestService()
	if _, err := service.Complete(context.Background(), "unknown", "code", ""); err == nil {
		t.Fatal("unknown state was accepted")
	}
	state, _ := startBrowserSSO(t, service)
	flows.flow.ExpiresAt = time.Now().Add(-time.Second)
	if _, err := service.Complete(context.Background(), state, "code", ""); err == nil {
		t.Fatal("expired state was accepted")
	}
	if central.exchangeCalls != 0 {
		t.Fatal("invalid state reached Central exchange")
	}

	flows.consumed = false
	flows.flow.ExpiresAt = time.Now().Add(time.Minute)
	result, err := service.Complete(context.Background(), state, "", "pending")
	if err != nil || result.Status != "pending" || result.Session != nil {
		t.Fatalf("pending result failed: %#v %v", result, err)
	}
	if _, err := service.Complete(context.Background(), state, "", "pending"); err == nil {
		t.Fatal("state replay was accepted")
	}
}

func TestBrowserSSOExchangeUsesTrustedContractAndLinkedIdentity(t *testing.T) {
	service, _, central, users := browserSSOTestService()
	state, _ := startBrowserSSO(t, service)
	result, err := service.Complete(context.Background(), state, "central-code", "")
	if err != nil {
		t.Fatal(err)
	}
	key, _ := service.key()
	if central.exchangeCalls != 1 || central.code != "central-code" ||
		central.clientID != "klinik-pkp-web" ||
		central.redirectURI != "http://localhost:3000/auth/sso/callback" ||
		central.codeVerifier != deriveKlinikPKCEVerifier(key, state) {
		t.Fatal("Central exchange did not use trusted client, redirect, and PKCE values")
	}
	if central.contextCalls != 1 || central.contextToken != "central-access" || users.requested != testAppUserID {
		t.Fatal("approved Klinik context was not resolved through app_user_id")
	}
	if result.Session.AccessToken != "central-access" || result.User.ID != testAppUserID ||
		result.User.AppUserID != testAppUserID || result.User.Role != "user" ||
		result.User.NIP != "198001012010011001" {
		t.Fatalf("unexpected authenticated result: %#v", result)
	}
}

func TestBrowserSSOMembershipAndIdentityFailures(t *testing.T) {
	for _, status := range []string{"denied", "rejected", "disabled"} {
		t.Run(status, func(t *testing.T) {
			service, _, central, _ := browserSSOTestService()
			state, _ := startBrowserSSO(t, service)
			_, err := service.Complete(context.Background(), state, "", status)
			var ssoError *BrowserSSOError
			if !errors.As(err, &ssoError) || ssoError.StatusCode != http.StatusForbidden || central.exchangeCalls != 0 {
				t.Fatalf("membership status was not denied safely: %v", err)
			}
		})
	}

	service, _, central, _ := browserSSOTestService()
	central.context.Role = "authenticated"
	state, _ := startBrowserSSO(t, service)
	if _, err := service.Complete(context.Background(), state, "code", ""); err == nil {
		t.Fatal("built-in Supabase role was accepted as Klinik authorization")
	}

	service, _, _, users := browserSSOTestService()
	users.canonical = nil
	state, _ = startBrowserSSO(t, service)
	_, err := service.Complete(context.Background(), state, "code", "")
	var ssoError *BrowserSSOError
	if !errors.As(err, &ssoError) || ssoError.Code != "KLINIK_ACCESS_DENIED" {
		t.Fatalf("unlinked identity did not return controlled error: %v", err)
	}
}

func TestBrowserSSOApprovedCanonicalIdentityNeedsNoLocalUser(t *testing.T) {
	service, _, _, users := browserSSOTestService()
	state, _ := startBrowserSSO(t, service)
	result, err := service.Complete(context.Background(), state, "code", "")
	if err != nil {
		t.Fatal(err)
	}
	if result.User.ID != testAppUserID || users.canonical.AppUserID != testAppUserID {
		t.Fatalf("browser SSO did not use canonical app user identity: %#v", result.User)
	}
}

func TestBrowserSSOLogoutAndRoutesUseTrustedTargets(t *testing.T) {
	service, _, _, _ := browserSSOTestService()
	handler := NewHandler(NewService(nil), service)
	app := fiber.New()
	api := app.Group("/api/v1")
	SetupRoutes(api, handler)

	logoutRequest := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/auth/sso/logout?return_to=https://evil.example/&client_id=evil",
		nil,
	)
	logoutResponse, err := app.Test(logoutRequest)
	if err != nil {
		t.Fatal(err)
	}
	location, err := url.Parse(logoutResponse.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if logoutResponse.StatusCode != http.StatusFound || location.Path != "/logout" ||
		location.Query().Get("client_id") != "klinik-pkp-web" ||
		location.Query().Get("post_logout_redirect_uri") != "http://localhost:3000/login" ||
		strings.Contains(location.String(), "evil.example") {
		t.Fatalf("logout bridge accepted browser override: %d %q", logoutResponse.StatusCode, location.String())
	}

	startResponse, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/v1/auth/sso/start", nil))
	if err != nil {
		t.Fatal(err)
	}
	if startResponse.StatusCode != http.StatusFound {
		t.Fatalf("SSO start route unavailable: %d", startResponse.StatusCode)
	}
}

func TestBrowserSSODisabledFailsClosedAndKeepsLegacyService(t *testing.T) {
	disabled := NewBrowserSSOService(config.Config{}, nil, nil, nil)
	if _, err := disabled.LogoutURL(); err == nil {
		t.Fatal("disabled browser SSO produced a logout redirect")
	}
	legacy := NewService(nil)
	handler := NewHandler(legacy, disabled)
	if handler.service != legacy {
		t.Fatal("browser SSO wiring replaced the legacy authentication service")
	}
}

func TestBrowserSSOExchangeHandlerReturnsStructuredSession(t *testing.T) {
	service, _, _, _ := browserSSOTestService()
	state, _ := startBrowserSSO(t, service)
	handler := NewHandler(NewService(nil), service)
	app := fiber.New()
	app.Post("/api/v1/auth/sso/exchange", handler.BrowserSSOExchangeHandler)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/sso/exchange",
		strings.NewReader(`{"code":"central-code","state":"`+state+`"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	var payload utils.Response
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	data, ok := payload.Data.(map[string]any)
	if response.StatusCode != http.StatusOK || !payload.Success || !ok ||
		data["status"] != "authenticated" || data["accessToken"] != "central-access" ||
		data["refreshToken"] != "central-refresh" {
		t.Fatalf("unexpected exchange response: status=%d payload=%#v", response.StatusCode, payload)
	}
}
