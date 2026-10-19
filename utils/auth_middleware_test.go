package utils

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthMiddlewareRejectsMissingBearer(t *testing.T) {
	t.Setenv("SSO_BASE_URL", "")

	cases := []struct {
		name   string
		header string
	}{
		{name: "missing"},
		{name: "wrong scheme", header: "Basic credentials"},
		{name: "empty bearer", header: "Bearer "},
		{name: "malformed bearer", header: "Bearer a b"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var reached atomic.Bool

			app := fiber.New()

			app.Get("/private",
				AuthMiddleware(),
				func(c *fiber.Ctx) error {
					reached.Store(true)
					return c.SendStatus(http.StatusOK)
				},
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/private",
				nil,
			)

			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusUnauthorized {
				t.Fatalf(
					"expected 401, got %d",
					resp.StatusCode,
				)
			}

			if reached.Load() {
				t.Fatal("protected handler was reached")
			}
		})
	}
}

func TestAuthMiddlewareDelegatesToCentralSSO(t *testing.T) {
	cases := []struct {
		name       string
		ssoStatus  int
		wantStatus int
	}{
		{
			name:       "invalid token",
			ssoStatus:  http.StatusUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "membership forbidden",
			ssoStatus:  http.StatusForbidden,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "central unavailable",
			ssoStatus:  http.StatusServiceUnavailable,
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "central internal error",
			ssoStatus:  http.StatusInternalServerError,
			wantStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			var reached atomic.Bool
			var requestValid atomic.Bool

			requestValid.Store(true)

			server := httptest.NewServer(
				http.HandlerFunc(func(
					w http.ResponseWriter,
					r *http.Request,
				) {
					calls.Add(1)

					if r.Method != http.MethodGet ||
						r.URL.Path !=
							"/api/v1/auth/klinik-context" ||
						r.Header.Get("Authorization") !=
							"Bearer legacy-style-token" {

						requestValid.Store(false)
					}

					w.WriteHeader(tc.ssoStatus)
				}),
			)
			defer server.Close()

			t.Setenv("SSO_BASE_URL", server.URL)

			app := fiber.New()

			app.Get(
				"/private",
				AuthMiddleware(),
				func(c *fiber.Ctx) error {
					reached.Store(true)
					return c.SendStatus(http.StatusOK)
				},
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/private",
				nil,
			)

			req.Header.Set(
				"Authorization",
				"Bearer legacy-style-token",
			)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf(
					"expected HTTP %d, got %d",
					tc.wantStatus,
					resp.StatusCode,
				)
			}

			if calls.Load() != 1 {
				t.Fatalf(
					"expected one SSO request, got %d",
					calls.Load(),
				)
			}

			if !requestValid.Load() {
				t.Fatal("unexpected SSO request")
			}

			if reached.Load() {
				t.Fatal(
					"protected handler was reached after SSO rejection",
				)
			}
		})
	}
}
