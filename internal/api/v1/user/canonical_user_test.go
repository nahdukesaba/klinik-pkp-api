package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type fakeDirectoryStore struct {
	users         []DirectoryUser
	query         DirectoryQuery
	getID         uuid.UUID
	updateActor   uuid.UUID
	updateTarget  uuid.UUID
	updateInput   DirectoryUpdate
	disableActor  uuid.UUID
	disableTarget uuid.UUID
	updateErr     error
	disableErr    error
}

func (f *fakeDirectoryStore) List(_ context.Context, query DirectoryQuery) ([]DirectoryUser, int64, error) {
	f.query = query
	return f.users, int64(len(f.users)), nil
}
func (f *fakeDirectoryStore) GetByAppUserID(_ context.Context, id uuid.UUID) (*DirectoryUser, error) {
	f.getID = id
	for i := range f.users {
		if f.users[i].AppUserID == id {
			result := f.users[i]
			return &result, nil
		}
	}
	return nil, nil
}
func (f *fakeDirectoryStore) Update(_ context.Context, actor, target uuid.UUID, input DirectoryUpdate) (*DirectoryUser, error) {
	f.updateActor, f.updateTarget, f.updateInput = actor, target, input
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	result := DirectoryUser{AppUserID: target, FullName: input.FullName, Role: input.Role, Status: input.Status}
	return &result, nil
}
func (f *fakeDirectoryStore) Disable(_ context.Context, actor, target uuid.UUID) error {
	f.disableActor, f.disableTarget = actor, target
	return f.disableErr
}

type fakeProvisioner struct {
	id    uuid.UUID
	token string
	input DirectoryCreate
	err   error
}

func (f *fakeProvisioner) Provision(_ context.Context, token string, input DirectoryCreate) (uuid.UUID, error) {
	f.token, f.input = token, input
	return f.id, f.err
}

func directoryFixture(status MembershipStatus) DirectoryUser {
	authID := uuid.MustParse("22222222-2222-4222-8222-222222222222")
	nip := "198001012010011001"
	return DirectoryUser{AppUserID: uuid.New(), AuthUserID: &authID, FullName: "Pegawai Klinik", Email: "pegawai@example.com", NIP: &nip, Role: RoleAdmin, Status: status, CreatedAt: time.Now(), UpdatedAt: time.Now()}
}

func TestCanonicalDirectoryIncludesEveryMembershipLifecycleStatus(t *testing.T) {
	store := &fakeDirectoryStore{users: []DirectoryUser{directoryFixture(StatusApproved), directoryFixture(StatusPending), directoryFixture(StatusRejected), directoryFixture(StatusDisabled)}}
	service := newServiceWithStores(nil, nil, store, nil)
	users, total, err := service.List(context.Background(), DirectoryQuery{Pagination: utils.PaginationQuery{Page: 1, Limit: 10}})
	if err != nil || total != 4 || len(users) != 4 {
		t.Fatalf("unexpected canonical list: total=%d users=%d err=%v", total, len(users), err)
	}
	for index, status := range []MembershipStatus{StatusApproved, StatusPending, StatusRejected, StatusDisabled} {
		if users[index].Status != status {
			t.Fatalf("status %s missing", status)
		}
	}
}

func TestCanonicalDirectoryFiltersAreMembershipFields(t *testing.T) {
	store := &fakeDirectoryStore{}
	service := newServiceWithStores(nil, nil, store, nil)
	_, _, err := service.List(context.Background(), DirectoryQuery{Pagination: utils.PaginationQuery{Page: 2, Limit: 20}, Keyword: " nip ", Role: RoleUser, Status: StatusPending})
	if err != nil {
		t.Fatal(err)
	}
	if store.query.Keyword != "nip" || store.query.Role != RoleUser || store.query.Status != StatusPending {
		t.Fatalf("filters not canonical: %#v", store.query)
	}
	if _, _, err := service.List(context.Background(), DirectoryQuery{Role: "superadmin"}); !errors.Is(err, ErrInvalidDirectoryMutation) {
		t.Fatalf("invalid role accepted: %v", err)
	}
}

func TestCanonicalDetailUsesAppUserIDAndEmployeeProfile(t *testing.T) {
	user := directoryFixture(StatusApproved)
	store := &fakeDirectoryStore{users: []DirectoryUser{user}}
	result, err := newServiceWithStores(nil, nil, store, nil).Get(context.Background(), user.AppUserID)
	if err != nil || result == nil || result.AppUserID != user.AppUserID || result.NIP == nil || *result.NIP != "198001012010011001" || result.Role != RoleAdmin {
		t.Fatalf("unexpected canonical detail: %#v err=%v", result, err)
	}
	legacyLocalID := uuid.MustParse("33333333-3333-4333-8333-333333333333")
	if _, err := newServiceWithStores(nil, nil, store, nil).Get(context.Background(), legacyLocalID); !errors.Is(err, ErrDirectoryUserNotFound) {
		t.Fatalf("legacy local ID unexpectedly resolved: %v", err)
	}
}

func TestCanonicalUpdateCarriesProfileAndMembershipOnly(t *testing.T) {
	actor, target := uuid.New(), uuid.New()
	store := &fakeDirectoryStore{}
	input := DirectoryUpdate{FullName: " Pegawai Baru ", Phone: " 0812 ", NIP: " 123 ", NIK: " 456 ", Alamat: " Jalan ", GelarDepan: " Dr. ", GelarBelakang: " M.T. ", Role: RoleUser, Status: StatusRejected}
	_, err := newServiceWithStores(nil, nil, store, nil).Update(context.Background(), actor, target, input)
	if err != nil {
		t.Fatal(err)
	}
	if store.updateActor != actor || store.updateTarget != target || store.updateInput.FullName != "Pegawai Baru" || store.updateInput.NIP != "123" || store.updateInput.Role != RoleUser || store.updateInput.Status != StatusRejected {
		t.Fatalf("canonical mutation changed: %#v", store.updateInput)
	}
}

func TestDisableTargetsOnlyCanonicalMembership(t *testing.T) {
	actor, target := uuid.New(), uuid.New()
	store := &fakeDirectoryStore{}
	if err := newServiceWithStores(nil, nil, store, nil).Disable(context.Background(), actor, target); err != nil {
		t.Fatal(err)
	}
	if store.disableActor != actor || store.disableTarget != target {
		t.Fatal("disable did not use canonical IDs")
	}
	store.disableErr = ErrSelfAccessMutation
	if err := newServiceWithStores(nil, nil, store, nil).Disable(context.Background(), actor, actor); !errors.Is(err, ErrSelfAccessMutation) {
		t.Fatalf("self protection lost: %v", err)
	}
}

func TestCreateDelegatesCanonicalProvisioningWithoutPassword(t *testing.T) {
	id := uuid.New()
	user := directoryFixture(StatusPending)
	user.AppUserID = id
	store := &fakeDirectoryStore{users: []DirectoryUser{user}}
	provisioner := &fakeProvisioner{id: id}
	result, err := newServiceWithStores(nil, nil, store, provisioner).Create(context.Background(), "Bearer central-token", DirectoryCreate{Email: " User@Example.com ", DirectoryUpdate: DirectoryUpdate{FullName: " New User ", Role: RoleUser, Status: StatusPending}})
	if err != nil || result == nil {
		t.Fatalf("create failed: %v", err)
	}
	if provisioner.token != "Bearer central-token" || provisioner.input.Email != "user@example.com" || provisioner.input.FullName != "New User" || store.getID != id {
		t.Fatalf("unexpected provisioning contract: %#v", provisioner)
	}
	encoded, _ := json.Marshal(provisioner.input)
	if strings.Contains(strings.ToLower(string(encoded)), "password") {
		t.Fatal("canonical provisioning payload contains password")
	}
}

func TestCentralProvisionerForwardsBearerAndCanonicalPayload(t *testing.T) {
	id := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/klinik/users" || r.Header.Get("Authorization") != "Bearer trusted" {
			t.Errorf("unexpected Central request: %s %s", r.URL.Path, r.Header.Get("Authorization"))
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if _, exists := body["password"]; exists {
			t.Error("password forwarded to Central")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"appUserId": id}})
	}))
	defer server.Close()
	got, err := NewCentralProvisioner(server.URL).Provision(context.Background(), "Bearer trusted", DirectoryCreate{Email: "user@example.com", DirectoryUpdate: DirectoryUpdate{FullName: "User", Role: RoleUser, Status: StatusPending}})
	if err != nil || got != id {
		t.Fatalf("unexpected Central response: %s %v", got, err)
	}
}

func TestCanonicalAdminRepositoryHasNoLegacyDirectoryDependency(t *testing.T) {
	source, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ToLower(string(source))
	if strings.Contains(text, "klinikpkp_user") || strings.Contains(text, "klinikpkp_authentication") {
		t.Fatal("canonical repository references legacy Klinik tables")
	}
	for _, required := range []string{"app_users", "app_memberships", "applications", "app_pegawai", "klinik_pkp"} {
		if !strings.Contains(text, required) {
			t.Fatalf("canonical repository missing %s", required)
		}
	}
	if strings.Contains(text, "update public.app_users set role") || strings.Contains(text, "update public.app_users set status") {
		t.Fatal("app-level authorization may be written to app_users")
	}
}

func TestCanonicalRouteRejectsNonUUIDInsteadOfFallingBack(t *testing.T) {
	handler := NewHandler(newServiceWithStores(nil, nil, &fakeDirectoryStore{}, nil))
	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByIdHandler)
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/users/legacy-email@example.com", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.StatusCode)
	}
}
