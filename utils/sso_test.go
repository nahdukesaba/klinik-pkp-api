package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSSOClientExchangeBrowserSSOContract(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/v1/auth/sso/exchange" {
			t.Fatalf("unexpected Central request: %s %s", request.Method, request.URL.Path)
		}
		var payload map[string]string
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["code"] != "central-code" || payload["clientId"] != "klinik-pkp-web" ||
			payload["redirectUri"] != "http://localhost:3000/auth/sso/callback" ||
			payload["codeVerifier"] != "derived-verifier" {
			t.Fatalf("unexpected Central exchange payload: %#v", payload)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{
            "accessToken":"access",
            "refreshToken":"refresh",
            "tokenType":"bearer",
            "expiresIn":3600,
            "expiresAt":9999999999,
            "user":{"id":"11111111-1111-4111-8111-111111111111","email":"person@example.com"}
        }`))
	}))
	defer server.Close()

	client := NewSSOClientWithBaseURL(server.URL)
	session, err := client.ExchangeBrowserSSO(
		context.Background(),
		"central-code",
		"klinik-pkp-web",
		"http://localhost:3000/auth/sso/callback",
		"derived-verifier",
	)
	if err != nil {
		t.Fatal(err)
	}
	if session.AccessToken != "access" || session.RefreshToken != "refresh" || session.User.Email != "person@example.com" {
		t.Fatalf("unexpected Central session response: %#v", session)
	}
}
