package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns CORS middleware with configuration from env
func CORS() fiber.Handler {
	allowOrigins := getCORSAllowOrigins()

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
	})
}

func getCORSAllowOrigins() string {
	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins == "" {
		return "*" // Default untuk development
	}
	// Convert comma-separated to fiber format
	parts := strings.Split(origins, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return strings.Join(parts, ",")
}

