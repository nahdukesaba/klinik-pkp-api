package config

import (
	"encoding/base64"
	"testing"
	"time"
)

func validBrowserSSOConfig() Config {
	return Config{
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

func TestValidateBrowserSSO(t *testing.T) {
	if err := validBrowserSSOConfig().ValidateBrowserSSO(); err != nil {
		t.Fatalf("valid browser SSO config rejected: %v", err)
	}
	disabled := Config{}
	if err := disabled.ValidateBrowserSSO(); err != nil {
		t.Fatalf("disabled browser SSO broke unrelated startup: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*Config)
	}{
		{"empty client", func(c *Config) { c.SSOClientID = "" }},
		{"relative base URL", func(c *Config) { c.SSOBaseURL = "/central" }},
		{"wildcard redirect", func(c *Config) { c.SSORedirectURI = "https://*.example.test/callback" }},
		{"invalid post logout scheme", func(c *Config) { c.SSOPostLogoutRedirectURI = "javascript:alert(1)" }},
		{"short TTL", func(c *Config) { c.SSOFlowTTL = time.Second }},
		{"long TTL", func(c *Config) { c.SSOFlowTTL = time.Hour }},
		{"malformed key", func(c *Config) { c.SSOFlowHMACKey = "not-base64" }},
		{"short key", func(c *Config) {
			c.SSOFlowHMACKey = base64.StdEncoding.EncodeToString([]byte("short"))
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := validBrowserSSOConfig()
			test.mutate(&cfg)
			if err := cfg.ValidateBrowserSSO(); err == nil {
				t.Fatal("invalid browser SSO config was accepted")
			}
		})
	}
}

func TestLoadConfigBrowserSSOValues(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("SSO_BROWSER_ENABLED", "true")
	t.Setenv("SSO_BASE_URL", " http://127.0.0.1:18082 ")
	t.Setenv("SSO_CLIENT_ID", " klinik-pkp-web ")
	t.Setenv("SSO_INTROSPECTION_SECRET", " klinik-test-introspection-secret-00001 ")
	t.Setenv("SSO_REDIRECT_URI", " http://localhost:3000/auth/sso/callback ")
	t.Setenv("SSO_POST_LOGOUT_REDIRECT_URI", " http://localhost:3000/login ")
	t.Setenv("SSO_FLOW_TTL", "5m")
	t.Setenv("SSO_FLOW_HMAC_KEY", validBrowserSSOConfig().SSOFlowHMACKey)

	cfg := LoadConfig()
	if !cfg.SSOBrowserEnabled || cfg.SSOBaseURL != "http://127.0.0.1:18082" ||
		cfg.SSOClientID != "klinik-pkp-web" ||
		cfg.SSORedirectURI != "http://localhost:3000/auth/sso/callback" ||
		cfg.SSOPostLogoutRedirectURI != "http://localhost:3000/login" ||
		cfg.SSOFlowTTL != 5*time.Minute {
		t.Fatal("browser SSO environment values were not loaded correctly")
	}
	if err := cfg.ValidateBrowserSSO(); err != nil {
		t.Fatalf("loaded browser SSO config is invalid: %v", err)
	}
}
