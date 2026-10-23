package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"klinik-pkp-api/internal/api/v1/user"
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	serviceLocalUserID = uuid.MustParse("cccccccc-cccc-4ccc-8ccc-cccccccccccc")
	serviceAppUserID   = uuid.MustParse("dddddddd-dddd-4ddd-8ddd-dddddddddddd")
)

type fakeLegacyAuthenticationStore struct {
	localUser       *user.User
	authentication  *Authentication
	findUserErr     error
	findSessionErr  error
	createErr       error
	rotateErr       error
	deleteErr       error
	createdToken    string
	rotatedOldToken string
	rotatedNewToken string
	deletedHash     string
	deletedPlain    string
}

func (f *fakeLegacyAuthenticationStore) FindUserByEmail(context.Context, string) (*user.User, error) {
	return f.localUser, f.findUserErr
}

func (f *fakeLegacyAuthenticationStore) FindAuthentication(
	context.Context,
	string,
	string,
) (*Authentication, error) {
	return f.authentication, f.findSessionErr
}

func (f *fakeLegacyAuthenticationStore) CreateAuthentication(_ context.Context, token string) error {
	f.createdToken = token
	return f.createErr
}

func (f *fakeLegacyAuthenticationStore) RotateAuthentication(
	_ context.Context,
	oldToken string,
	newToken string,
) error {
	f.rotatedOldToken = oldToken
	f.rotatedNewToken = newToken
	return f.rotateErr
}

func (f *fakeLegacyAuthenticationStore) DeleteAuthentication(
	_ context.Context,
	hashedToken string,
	plainToken string,
) error {
	f.deletedHash = hashedToken
	f.deletedPlain = plainToken
	return f.deleteErr
}

type fakeServiceAuthorizationResolver struct {
	trusted *utils.AuthContext
	err     error
	calls   int
}

func (f *fakeServiceAuthorizationResolver) ResolveKlinikAuthorization(
	context.Context,
	uuid.UUID,
) (*utils.AuthContext, error) {
	f.calls++
	if f.trusted == nil {
		return nil, f.err
	}
	copyContext := *f.trusted
	return &copyContext, f.err
}

func (f *fakeServiceAuthorizationResolver) ResolveLegacyKlinikAuthorization(
	context.Context,
	uuid.UUID,
) (*utils.AuthContext, error) {
	f.calls++
	if f.trusted == nil {
		return nil, f.err
	}
	copyContext := *f.trusted
	return &copyContext, f.err
}

func setLegacyJWTTestConfig(t *testing.T) {
	t.Helper()
	t.Setenv("ACCESS_TOKEN_KEY", "test-access-token-key-32-characters-minimum")
	t.Setenv("REFRESH_TOKEN_KEY", "test-refresh-token-key-32-characters-minimum")
	t.Setenv("ACCESS_TOKEN_EXPIRATION", "15m")
	t.Setenv("REFRESH_TOKEN_EXPIRATION", "24h")
}

func legacyTestUser(t *testing.T, localRole string) *user.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	linkedID := serviceAppUserID
	return &user.User{
		ID:        serviceLocalUserID,
		AppUserID: &linkedID,
		Email:     "person@example.com",
		NIP:       "19800101",
		Password:  string(hash),
		Role:      localRole,
		IsActive:  true,
	}
}

func serviceAuthorization(role string) *utils.AuthContext {
	localID := serviceLocalUserID
	return &utils.AuthContext{
		LocalUserID:      &localID,
		AppUserID:        serviceAppUserID,
		Email:            "person@example.com",
		Role:             role,
		MembershipStatus: "approved",
	}
}

func TestLegacyLoginRequiresApprovedMembershipAfterBcrypt(t *testing.T) {
	setLegacyJWTTestConfig(t)
	store := &fakeLegacyAuthenticationStore{localUser: legacyTestUser(t, "admin")}
	authorization := &fakeServiceAuthorizationResolver{err: utils.ErrKlinikAuthorizationDenied}
	service := newServiceWithDependencies(nil, store, authorization)

	_, err := service.AddAuthentication(context.Background(), AddAuthenticationPayload{
		Email: "person@example.com", NIP: "19800101", Password: "correct-password",
	})
	if !errors.Is(err, utils.ErrKlinikAuthorizationDenied) || authorization.calls != 1 || store.createdToken != "" {
		t.Fatalf("rejected membership was not denied after password validation: %v", err)
	}
}

func TestLegacyLoginIssuesMembershipRoleTokens(t *testing.T) {
	setLegacyJWTTestConfig(t)
	store := &fakeLegacyAuthenticationStore{localUser: legacyTestUser(t, "admin")}
	service := newServiceWithDependencies(nil, store, &fakeServiceAuthorizationResolver{
		trusted: serviceAuthorization("user"),
	})

	result, err := service.AddAuthentication(context.Background(), AddAuthenticationPayload{
		Email: "person@example.com", NIP: "19800101", Password: "correct-password",
	})
	if err != nil {
		t.Fatal(err)
	}
	accessToken, accessOK := result["access_token"].(string)
	refreshToken, refreshOK := result["refresh_token"].(string)
	claims, verifyErr := utils.NewTokenManager().VerifyAccessToken(accessToken)
	if !accessOK || !refreshOK || verifyErr != nil || claims.Role != "user" ||
		store.createdToken != hashRefreshToken(refreshToken) || store.createdToken == refreshToken {
		t.Fatalf("legacy login did not issue membership-derived tokens: result=%#v claims=%#v err=%v", result, claims, verifyErr)
	}
}

func TestLegacyRefreshRevalidatesMembership(t *testing.T) {
	setLegacyJWTTestConfig(t)
	refreshToken, err := utils.NewTokenManager().GenerateRefreshToken(
		serviceLocalUserID,
		"person@example.com",
		"user",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("rejected membership denies refresh", func(t *testing.T) {
		store := &fakeLegacyAuthenticationStore{authentication: &Authentication{Token: hashRefreshToken(refreshToken)}}
		service := newServiceWithDependencies(nil, store, &fakeServiceAuthorizationResolver{
			err: utils.ErrKlinikAuthorizationDenied,
		})
		if _, err := service.UpdateAuthentication(context.Background(), refreshToken); !errors.Is(err, utils.ErrKlinikAuthorizationDenied) {
			t.Fatalf("refresh was not denied: %v", err)
		}
		if store.rotatedNewToken != "" {
			t.Fatal("denied refresh rotated the session")
		}
	})

	t.Run("approved refresh uses current membership role", func(t *testing.T) {
		store := &fakeLegacyAuthenticationStore{authentication: &Authentication{Token: hashRefreshToken(refreshToken)}}
		service := newServiceWithDependencies(nil, store, &fakeServiceAuthorizationResolver{
			trusted: serviceAuthorization("admin"),
		})
		result, err := service.UpdateAuthentication(context.Background(), refreshToken)
		if err != nil {
			t.Fatal(err)
		}
		accessClaims, err := utils.NewTokenManager().VerifyAccessToken(result["access_token"].(string))
		if err != nil || accessClaims.Role != "admin" || store.rotatedOldToken != hashRefreshToken(refreshToken) ||
			store.rotatedNewToken != hashRefreshToken(result["refresh_token"].(string)) {
			t.Fatalf("refresh kept stale authorization: claims=%#v old=%q new=%q err=%v", accessClaims, store.rotatedOldToken, store.rotatedNewToken, err)
		}
	})
}

func TestContextualLogoutContract(t *testing.T) {
	tests := []struct {
		name        string
		source      utils.LoginSource
		query       string
		wantCentral bool
	}{
		{name: "direct login", source: utils.LoginSourceDirect, wantCentral: false},
		{name: "browser SSO login", source: utils.LoginSourceSSO, wantCentral: true},
		{name: "browser cannot override direct source", source: utils.LoginSourceDirect, query: "?login_source=sso&centralLogoutRequired=true", wantCentral: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := NewHandler(newServiceWithDependencies(
				nil,
				&fakeLegacyAuthenticationStore{},
				&fakeServiceAuthorizationResolver{},
			))
			app := fiber.New()
			app.Delete("/logout", func(ctx *fiber.Ctx) error {
				ctx.Locals("auth_context", &utils.AuthContext{LoginSource: test.source})
				return ctx.Next()
			}, handler.DeleteAuthenticationHandler)
			response, err := app.Test(httptest.NewRequest(http.MethodDelete, "/logout"+test.query, nil))
			if err != nil {
				t.Fatal(err)
			}
			var payload struct {
				Success bool `json:"success"`
				Data    struct {
					CentralLogoutRequired bool `json:"centralLogoutRequired"`
				} `json:"data"`
			}
			if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != http.StatusOK || !payload.Success || payload.Data.CentralLogoutRequired != test.wantCentral {
				t.Fatalf("unexpected logout response: status=%d payload=%#v", response.StatusCode, payload)
			}
		})
	}
}
