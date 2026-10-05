package authentication

import (
	"klinik-pkp-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SETUP ROUTES CONFIGURES AUTHENTICATION ENDPOINTS
func SetupRoutes(app fiber.Router, handler *Handler) {
	auth := app.Group("/authentications")

	loginLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			if ip := ctx.Get("X-Real-IP"); ip != "" {
				return ip
			}
			return ctx.IP()
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return utils.JSONResponse(
				ctx,
				fiber.StatusTooManyRequests,
				"too many login attempts, please try again later",
				fiber.ErrTooManyRequests,
				true,
			)
		},
	})

	auth.Post("/", loginLimiter, handler.PostAuthenticationHandler)
	auth.Put("/", handler.PutAuthenticationHandler)
	auth.Delete("/", handler.DeleteAuthenticationHandler)
}
