package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// APP
	AppEnv       string
	DockerTarget string
	DockerCmd    string
	AppPort      string

	// DATABASE
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	DBSSLMode       string
	SUPABASE_DB_URL string

	// ROLES
	RoleSuperAdmin         string
	RoleAdminEselon1       string
	RoleVerificatorEselon1 string
	RoleAdminBalai         string
	RoleVerificatorBalai   string
	RoleSurveyor           string
	RoleUser               string

	// RESOURCES
	ResourceNegara       string
	ResourcePengembang   string
	ResourceSwadaya      string
	ResourceGotongroyong string

	// BANNED WORDS
	BannedWords []string

	// CORS
	CORSAllowOrigins []string

	// CENTRAL BROWSER SSO
	SSOBrowserEnabled        bool
	SSOBaseURL               string
	SSOClientID              string
	SSORedirectURI           string
	SSOPostLogoutRedirectURI string
	SSOFlowTTL               time.Duration
	SSOFlowHMACKey           string
}

func LoadConfig() Config {
	godotenv.Load()

	cfg := Config{
		// APP
		AppEnv:       getEnv("APP_ENV", "development"),
		DockerTarget: getEnv("DOCKER_TARGET", "dev"),
		DockerCmd:    getEnv("DOCKER_CMD", "air"),
		AppPort:      getEnv("APP_PORT", "8000"),

		// DATABASE
		DBHost:    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", "klinik-pkp-api"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),

		// ROLES
		RoleSuperAdmin:         getEnv("ROLE_SUPER_ADMIN", "Super Admin"),
		RoleAdminEselon1:       getEnv("ROLE_ADMIN_ESELON_1", "Admin Eselon 1"),
		RoleVerificatorEselon1: getEnv("ROLE_VERIFICATOR_ESELON_1", "Verificator Eselon 1"),
		RoleAdminBalai:         getEnv("ROLE_ADMIN_BALAI", "Admin Balai"),
		RoleVerificatorBalai:   getEnv("ROLE_VERIFICATOR_BALAI", "Verificator Balai"),
		RoleSurveyor:           getEnv("ROLE_SURVEYOR", "Surveyor"),
		RoleUser:               getEnv("ROLE_USER", "User"),

		// RESOURCES
		ResourceNegara:       getEnv("RESOURCE_NEGARA", "Negara"),
		ResourcePengembang:   getEnv("RESOURCE_PENGEMBANG", "Pengembang"),
		ResourceSwadaya:      getEnv("RESOURCE_SWADAYA", "Swadaya"),
		ResourceGotongroyong: getEnv("RESOURCE_GOTONGROYONG", "Gotong Royong"),

		// BANNED WORDS
		BannedWords: parseCSV(getEnv("BANNED_WORDS", "")),

		// CORS
		CORSAllowOrigins: parseCSV(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173")),

		// CENTRAL BROWSER SSO
		SSOBrowserEnabled:        getEnvBool("SSO_BROWSER_ENABLED", false),
		SSOBaseURL:               strings.TrimSpace(getEnv("SSO_BASE_URL", "")),
		SSOClientID:              strings.TrimSpace(getEnv("SSO_CLIENT_ID", "klinik-pkp-web")),
		SSORedirectURI:           strings.TrimSpace(getEnv("SSO_REDIRECT_URI", "")),
		SSOPostLogoutRedirectURI: strings.TrimSpace(getEnv("SSO_POST_LOGOUT_REDIRECT_URI", "")),
		SSOFlowTTL:               getEnvDuration("SSO_FLOW_TTL", 5*time.Minute),
		SSOFlowHMACKey:           strings.TrimSpace(getEnv("SSO_FLOW_HMAC_KEY", "")),
	}

	// VALIDATE REQUIRED CONFIGURATIONS
	if cfg.DBHost == "" || cfg.DBName == "" {
		log.Fatal("Database configuration is required in .env file")
	}

	return cfg
}

func (c Config) ValidateBrowserSSO() error {
	if !c.SSOBrowserEnabled {
		return nil
	}
	if strings.TrimSpace(c.SSOClientID) == "" {
		return fmt.Errorf("SSO_CLIENT_ID is required when browser SSO is enabled")
	}
	for name, value := range map[string]string{
		"SSO_BASE_URL":                 c.SSOBaseURL,
		"SSO_REDIRECT_URI":             c.SSORedirectURI,
		"SSO_POST_LOGOUT_REDIRECT_URI": c.SSOPostLogoutRedirectURI,
	} {
		if !validAbsoluteHTTPURL(value) {
			return fmt.Errorf("%s must be an absolute HTTP(S) URL", name)
		}
	}
	if c.SSOFlowTTL < time.Minute || c.SSOFlowTTL > 15*time.Minute {
		return fmt.Errorf("SSO_FLOW_TTL must be between 1m and 15m")
	}
	key, err := base64.StdEncoding.DecodeString(c.SSOFlowHMACKey)
	if err != nil {
		key, err = base64.RawStdEncoding.DecodeString(c.SSOFlowHMACKey)
	}
	if err != nil || len(key) < 32 {
		return fmt.Errorf("SSO_FLOW_HMAC_KEY must be a base64-encoded key of at least 32 bytes")
	}
	return nil
}

func validAbsoluteHTTPURL(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || strings.ContainsAny(value, "*#") {
		return false
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Host == "" || parsed.User != nil {
		return false
	}
	scheme := strings.ToLower(parsed.Scheme)
	return scheme == "http" || scheme == "https"
}

// GET DATA SOURCE NAME FOR DATABASE CONNECTION
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSSLMode)
}

// GET ENVIRONMENT VARIABLE OR RETURN DEFAULT VALUE
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return parsed
}

// SPLIT A COMMA-SEPARATED STRING INTO A SLICE OF STRINGS
func parseCSV(value string) []string {
	if value == "" {
		return []string{}
	}
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
