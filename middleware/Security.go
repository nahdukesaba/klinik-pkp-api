package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent XSS attacks
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy
		c.Set("Content-Security-Policy", "default-src 'self'")
		
		// Strict Transport Security (for HTTPS)
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		// Remove X-Powered-By header
		c.Set("X-Powered-By", "")
		
		return c.Next()
	}
}
