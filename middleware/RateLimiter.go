package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// RateLimiter creates a rate limiting middleware
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,              // Max requests
		Expiration: 1 * time.Minute,  // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// AuthRateLimiter creates a stricter rate limiter for auth endpoints
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,               // Max 10 login attempts
		Expiration: 5 * time.Minute,  // Per 5 minutes
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many login attempts. Please try again in 5 minutes.",
			})
		},
	})
}
