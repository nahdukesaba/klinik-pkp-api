package utils

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"klinik-pkp-api/internal/identity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	middlewareLocalUserID = uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	middlewareAppUserID   = uuid.MustParse("bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb")
)

type fakeLocalTokenVerifier struct {
	claims *Claims
	err    error
}

func (f *fakeLocalTokenVerifier) VerifyAccessToken(string) (*Claims, error) {
	return f.claims, f.err
}

type fakeKlinikSSOResolver struct {
	resolved *SSOContext
	err      error
	calls    int
}

func (f *fakeKlinikSSOResolver) ResolveKlinikContext(context.Context, string) (*SSOContext, error) {
	f.calls++
	return f.resolved, f.err
}

type fakeKlinikAuthorizationResolver struct {
	context        *AuthContext
	err            error
	canonicalCalls int
	legacyCalls    int
}

type fakeCanonicalIdentityStore struct {
	canonical *identity.KlinikIdentity
	err       error
}

func (f *fakeCanonicalIdentityStore) GetByAppUserID(
	context.Context,
	uuid.UUID,
) (*identity.KlinikIdentity, error) {
	return f.canonical, f.err
}

func (f *fakeKlinikAuthorizationResolver) ResolveKlinikAuthorization(
	context.Context,
	uuid.UUID,
) (*AuthContext, error) {
	f.canonicalCalls++
	if f.context == nil {
		return nil, f.err
	}
	copyContext := *f.context
	return &copyContext, f.err
}

func (f *fakeKlinikAuthorizationResolver) ResolveLegacyKlinikAuthorization(
	context.Context,
	uuid.UUID,
) (*AuthContext, error) {
	f.legacyCalls++
	if f.context == nil {
		return nil, f.err
	}
	copyContext := *f.context
	return &copyContext, f.err
}

func approvedTestAuthContext(role string) *AuthContext {
	localID := middlewareLocalUserID
	return &AuthContext{
		LocalUserID:      &localID,
		AppUserID:        middlewareAppUserID,
		Email:            "person@example.com",
		Role:             role,
		MembershipStatus: "approved",
	}
}

func runProtectedRequest(
	t *testing.T,
	middleware fiber.Handler,
	after ...fiber.Handler,
) (int, string) {
	t.Helper()
	app := fiber.New()
	handlers := []fiber.Handler{middleware}
	handlers = append(handlers, after...)
	handlers = append(handlers, func(ctx *fiber.Ctx) error {
		trusted, ok := TrustedAuthContext(ctx)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
		return ctx.SendString(fmt.Sprintf("%s:%s", trusted.Role, trusted.LoginSource))
	})
	app.Get("/protected", handlers...)
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	body := make([]byte, response.ContentLength)
	_, _ = response.Body.Read(body)
	return response.StatusCode, string(body)
}

func TestAuthorizedKlinikContextMembershipRules(t *testing.T) {
	tests := []struct {
		name      string
		appUserID uuid.UUID
		status    string
		role      string
		wantRole  string
		wantErr   bool
	}{
		{name: "approved user", appUserID: middlewareAppUserID, status: "approved", role: "user", wantRole: "user"},
		{name: "approved admin", appUserID: middlewareAppUserID, status: "approved", role: "admin", wantRole: "admin"},
		{name: "pending", appUserID: middlewareAppUserID, status: "pending", role: "user", wantErr: true},
		{name: "rejected", appUserID: middlewareAppUserID, status: "rejected", role: "user", wantErr: true},
		{name: "disabled", appUserID: middlewareAppUserID, status: "disabled", role: "user", wantErr: true},
		{name: "missing membership", appUserID: uuid.Nil, status: "", role: "", wantErr: true},
		{name: "missing app user id", appUserID: uuid.Nil, status: "approved", role: "user", wantErr: true},
		{name: "invalid role", appUserID: middlewareAppUserID, status: "approved", role: "superadmin", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authUserID := uuid.MustParse("cccccccc-cccc-4ccc-8ccc-cccccccccccc")
			trusted, err := authorizedKlinikContext(&identity.KlinikIdentity{
				AppUserID: test.appUserID, AuthUserID: &authUserID,
				Email: "person@example.com", Role: test.role,
				MembershipStatus: test.status,
				ApplicationCode:  identity.KlinikApplicationCode,
			})
			if test.wantErr {
				if !errors.Is(err, ErrKlinikAuthorizationDenied) {
					t.Fatalf("expected authorization denial, got %v", err)
				}
				return
			}
			if err != nil || trusted.Role != test.wantRole {
				t.Fatalf("unexpected authorization result: %#v %v", trusted, err)
			}
		})
	}
}

func TestLocalJWTUsesAuthoritativeMembershipRole(t *testing.T) {
	tests := []struct {
		name           string
		claimRole      string
		membershipRole string
	}{
		{name: "local admin claim becomes membership user", claimRole: "admin", membershipRole: "user"},
		{name: "local user claim becomes membership admin", claimRole: "user", membershipRole: "admin"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			local := &fakeLocalTokenVerifier{claims: &Claims{UserID: middlewareLocalUserID, Role: test.claimRole}}
			central := &fakeKlinikSSOResolver{err: ErrSSOUnauthorized}
			authorization := &fakeKlinikAuthorizationResolver{context: approvedTestAuthContext(test.membershipRole)}
			status, body := runProtectedRequest(t, newAuthMiddleware(local, central, authorization, false))
			if status != http.StatusOK || body != test.membershipRole+":direct" || central.calls != 0 {
				t.Fatalf("role was not authoritative: status=%d body=%q centralCalls=%d", status, body, central.calls)
			}
		})
	}
}

func TestLocalJWTDeniesCurrentMembershipFailures(t *testing.T) {
	for _, name := range []string{"pending", "rejected", "disabled", "missing membership", "missing app_user_id", "invalid role"} {
		t.Run(name, func(t *testing.T) {
			status, _ := runProtectedRequest(t, newAuthMiddleware(
				&fakeLocalTokenVerifier{claims: &Claims{UserID: middlewareLocalUserID, Role: "admin"}},
				&fakeKlinikSSOResolver{},
				&fakeKlinikAuthorizationResolver{err: ErrKlinikAuthorizationDenied},
				false,
			))
			if status != http.StatusForbidden {
				t.Fatalf("expected forbidden, got %d", status)
			}
		})
	}
}

func TestExistingLocalJWTIsDeniedAfterMembershipRevocation(t *testing.T) {
	authorization := &fakeKlinikAuthorizationResolver{context: approvedTestAuthContext("user")}
	middleware := newAuthMiddleware(
		&fakeLocalTokenVerifier{claims: &Claims{UserID: middlewareLocalUserID}},
		&fakeKlinikSSOResolver{},
		authorization,
		false,
	)
	status, _ := runProtectedRequest(t, middleware)
	if status != http.StatusOK {
		t.Fatalf("approved membership was denied: %d", status)
	}
	authorization.context = nil
	authorization.err = ErrKlinikAuthorizationDenied
	status, _ = runProtectedRequest(t, middleware)
	if status != http.StatusForbidden || authorization.legacyCalls != 2 {
		t.Fatalf("existing token bypassed revalidation: status=%d calls=%d", status, authorization.legacyCalls)
	}
}

func TestCentralContextPathAndVerifiedAMRLoginSource(t *testing.T) {
	methods := []struct {
		name string
		amr  string
		want LoginSource
	}{
		{name: "password", amr: `[{"method":"password"}]`, want: LoginSourceDirect},
		{name: "magiclink", amr: `[{"method":"magiclink"}]`, want: LoginSourceSSO},
		{name: "otp", amr: `[{"method":"otp"}]`, want: LoginSourceSSO},
		{name: "sso", amr: `[{"method":"sso"}]`, want: LoginSourceSSO},
		{name: "saml", amr: `["saml"]`, want: LoginSourceSSO},
		{name: "unknown", amr: `[{"method":"webauthn"}]`, want: LoginSourceDirect},
		{name: "missing", amr: `null`, want: LoginSourceDirect},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			token := "header." + base64.RawURLEncoding.EncodeToString([]byte(`{"amr":`+method.amr+`}`)) + ".signature"
			middleware := &authMiddleware{
				local: &fakeLocalTokenVerifier{err: errors.New("not local")},
				central: &fakeKlinikSSOResolver{resolved: &SSOContext{
					ID: middlewareAppUserID, AppUserID: middlewareAppUserID,
					Email: "person@example.com", Role: "user", Status: "approved",
				}},
				authorization: &fakeKlinikAuthorizationResolver{
					context: approvedTestAuthContext("user"),
				},
			}
			trusted, err := middleware.authenticate(context.Background(), token)
			if err != nil || trusted.LoginSource != method.want {
				t.Fatalf("unexpected login source: %#v %v", trusted, err)
			}
		})
	}
}

func TestCentralTokenApprovedMembershipDoesNotRequireLocalUser(t *testing.T) {
	trusted := approvedTestAuthContext("user")
	trusted.LocalUserID = nil
	authorization := &fakeKlinikAuthorizationResolver{context: trusted}
	middleware := &authMiddleware{
		local: &fakeLocalTokenVerifier{err: errors.New("not local")},
		central: &fakeKlinikSSOResolver{resolved: &SSOContext{
			ID: middlewareAppUserID, AppUserID: middlewareAppUserID,
			Email: "person@example.com", Role: "user", Status: "approved",
		}},
		authorization: authorization,
	}

	result, err := middleware.authenticate(context.Background(), "central-token")
	if err != nil || result.AppUserID != middlewareAppUserID || result.LocalUserID != nil ||
		authorization.canonicalCalls != 1 || authorization.legacyCalls != 0 {
		t.Fatalf("canonical token required local identity: result=%#v err=%v", result, err)
	}
}

func TestDatabaseResolverCanonicalPathSucceedsWithoutLocalDatabase(t *testing.T) {
	authUserID := uuid.MustParse("cccccccc-cccc-4ccc-8ccc-cccccccccccc")
	phone := "08123456789"
	nip := "198001012010011001"
	resolver := &databaseKlinikAuthorizationResolver{
		identities: &fakeCanonicalIdentityStore{canonical: &identity.KlinikIdentity{
			AppUserID: middlewareAppUserID, AuthUserID: &authUserID,
			Email: "person@example.com", FullName: "Pegawai Klinik",
			Phone: &phone, NIP: &nip, Role: "admin",
			MembershipStatus: "approved",
			ApplicationCode:  identity.KlinikApplicationCode,
		}},
	}

	trusted, err := resolver.ResolveKlinikAuthorization(context.Background(), middlewareAppUserID)
	if err != nil || trusted.AppUserID != middlewareAppUserID || trusted.LocalUserID != nil ||
		trusted.Role != "admin" || trusted.NIP == nil || *trusted.NIP != nip {
		t.Fatalf("canonical resolver required local database: trusted=%#v err=%v", trusted, err)
	}
}

func TestCentralTokenRevokedMembershipFailsNextRequest(t *testing.T) {
	trusted := approvedTestAuthContext("user")
	trusted.LocalUserID = nil
	authorization := &fakeKlinikAuthorizationResolver{context: trusted}
	middleware := newAuthMiddleware(
		&fakeLocalTokenVerifier{err: errors.New("not local")},
		&fakeKlinikSSOResolver{resolved: &SSOContext{
			ID: middlewareAppUserID, AppUserID: middlewareAppUserID,
			Email: "person@example.com", Role: "user", Status: "approved",
		}},
		authorization,
		false,
	)
	status, _ := runProtectedRequest(t, middleware)
	if status != http.StatusOK {
		t.Fatalf("approved canonical membership was denied: %d", status)
	}
	authorization.context = nil
	authorization.err = ErrKlinikAuthorizationDenied
	status, _ = runProtectedRequest(t, middleware)
	if status != http.StatusForbidden || authorization.canonicalCalls != 2 {
		t.Fatalf("revoked canonical membership bypassed revalidation: status=%d calls=%d", status, authorization.canonicalCalls)
	}
}

func TestAdminRouteUsesMembershipRole(t *testing.T) {
	tests := []struct {
		name           string
		claimRole      string
		membershipRole string
		wantStatus     int
	}{
		{name: "admin claim cannot elevate user membership", claimRole: "admin", membershipRole: "user", wantStatus: http.StatusForbidden},
		{name: "user claim receives current admin membership", claimRole: "user", membershipRole: "admin", wantStatus: http.StatusOK},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			middleware := newAuthMiddleware(
				&fakeLocalTokenVerifier{claims: &Claims{UserID: middlewareLocalUserID, Role: test.claimRole}},
				&fakeKlinikSSOResolver{},
				&fakeKlinikAuthorizationResolver{context: approvedTestAuthContext(test.membershipRole)},
				false,
			)
			status, _ := runProtectedRequest(t, middleware, RoleMiddleware("admin"))
			if status != test.wantStatus {
				t.Fatalf("expected status %d, got %d", test.wantStatus, status)
			}
		})
	}
}
