package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"klinik-pkp-api/internal/identity"
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type fakeCurrentUserIdentities struct {
	canonical *identity.KlinikIdentity
	err       error
}

func (f *fakeCurrentUserIdentities) GetByAppUserID(
	context.Context,
	uuid.UUID,
) (*identity.KlinikIdentity, error) {
	return f.canonical, f.err
}

func TestUsersMeReturnsCanonicalDataWithoutLocalUser(t *testing.T) {
	appUserID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	authUserID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	phone := "08123456789"
	nip := "198001012010011001"
	store := &fakeCurrentUserIdentities{canonical: &identity.KlinikIdentity{
		AppUserID: appUserID, AuthUserID: &authUserID,
		Email: "person@example.com", FullName: "Pegawai Klinik",
		Phone: &phone, NIP: &nip, Role: "admin", MembershipStatus: "approved",
		ApplicationCode: identity.KlinikApplicationCode,
	}}
	handler := NewHandler(NewServiceWithIdentityStore(nil, store))
	app := fiber.New()
	app.Get("/users/:id", func(ctx *fiber.Ctx) error {
		ctx.Locals("auth_context", &utils.AuthContext{
			AppUserID: appUserID, AuthUserID: &authUserID,
			Role: "admin", MembershipStatus: "approved",
		})
		return ctx.Next()
	}, handler.GetUserByIdHandler)

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/users/me", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var payload struct {
		Success bool        `json:"success"`
		Data    CurrentUser `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || !payload.Success ||
		payload.Data.ID != appUserID || payload.Data.AppUserID != appUserID ||
		payload.Data.AuthUserID != authUserID || payload.Data.Name != "Pegawai Klinik" ||
		payload.Data.NIP != nip || payload.Data.Role != "admin" || !payload.Data.IsActive {
		t.Fatalf("unexpected canonical /users/me response: status=%d payload=%#v", response.StatusCode, payload)
	}
}

func TestUsersMeDeniesNonApprovedCanonicalIdentity(t *testing.T) {
	appUserID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	authUserID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	for _, status := range []string{"pending", "rejected", "disabled"} {
		t.Run(status, func(t *testing.T) {
			store := &fakeCurrentUserIdentities{canonical: &identity.KlinikIdentity{
				AppUserID: appUserID, AuthUserID: &authUserID,
				Role: "user", MembershipStatus: status,
				ApplicationCode: identity.KlinikApplicationCode,
			}}
			_, err := NewServiceWithIdentityStore(nil, store).GetCurrentUser(context.Background(), appUserID)
			if err != ErrCanonicalIdentityUnavailable {
				t.Fatalf("expected canonical identity unavailable, got %v", err)
			}
		})
	}
}
