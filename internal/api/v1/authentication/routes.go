package authentication

import (
	"klinik-pkp-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// SETUP ROUTES CONFIGURES AUTHENTICATION ENDPOINTS
func SetupRoutes(app fiber.Router, handler *Handler) {
	authentications := app.Group("/authentications")

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

	authentications.Post("/", loginLimiter, handler.PostAuthenticationHandler)
	authentications.Put("/", handler.PutAuthenticationHandler)
	authentications.Delete("/", handler.DeleteAuthenticationHandler)

	browserSSO := app.Group("/auth/sso")
	browserSSO.Get("/start", handler.BrowserSSOStartHandler)
	browserSSO.Post("/exchange", loginLimiter, handler.BrowserSSOExchangeHandler)
	browserSSO.Get("/logout", handler.BrowserSSOLogoutHandler)
}
