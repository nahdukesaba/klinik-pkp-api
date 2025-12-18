package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// App
	AppEnv       string
	DockerTarget string
	DockerCmd    string

	// Database
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// Seed
	DBSeed bool

	// Roles (Database roles only)
	RoleSuperAdmin          string
	RoleAdminEselon1        string
	RoleVerificatorEselon1  string
	RoleAdminBalai          string
	RoleVerificatorBalai    string
	RoleSurveyor            string
	RoleUser                string

	// Resources
	ResourceNegara       string
	ResourcePengembang   string
	ResourceSwadaya      string
	ResourceGotongroyong string

	// Security
	BannedWords []string

	// CORS
	CORSAllowOrigins []string
}

func LoadConfig() Config {
	godotenv.Load()

	cfg := Config{
		// App
		AppEnv:       getEnv("APP_ENV", "development"),
		DockerTarget: getEnv("DOCKER_TARGET", "dev"),
		DockerCmd:    getEnv("DOCKER_CMD", "air"),

		// Database
		DBHost: getEnv("DB_HOST", "127.0.0.1"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", ""),
		DBName: getEnv("DB_NAME", "klinik-api"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "supersecretjwtkey"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		// Seed
		DBSeed: getEnv("DB_SEED", "true") == "true",

		// Roles (Database roles only)
		RoleSuperAdmin:         getEnv("ROLE_SUPER_ADMIN", "Super Admin"),
		RoleAdminEselon1:       getEnv("ROLE_ADMIN_ESELON_1", "Admin Eselon 1"),
		RoleVerificatorEselon1: getEnv("ROLE_VERIFICATOR_ESELON_1", "Verificator Eselon 1"),
		RoleAdminBalai:         getEnv("ROLE_ADMIN_BALAI", "Admin Balai"),
		RoleVerificatorBalai:   getEnv("ROLE_VERIFICATOR_BALAI", "Verificator Balai"),
		RoleSurveyor:           getEnv("ROLE_SURVEYOR", "Surveyor"),
		RoleUser:               getEnv("ROLE_USER", "User"),

		// Resources
		ResourceNegara:       getEnv("RESOURCE_NEGARA", "Negara"),
		ResourcePengembang:   getEnv("RESOURCE_PENGEMBANG", "Pengembang"),
		ResourceSwadaya:      getEnv("RESOURCE_SWADAYA", "Swadaya"),
		ResourceGotongroyong: getEnv("RESOURCE_GOTONGROYONG", "Gotong Royong"),

		// Security
		BannedWords: parseCSV(getEnv("BANNED_WORDS", "")),

		// CORS
		CORSAllowOrigins: parseCSV(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5173")),
	}

	// Validate required fields
	if cfg.DBHost == "" || cfg.DBName == "" {
		log.Fatal("Database configuration is required in .env file")
	}

	return cfg
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

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