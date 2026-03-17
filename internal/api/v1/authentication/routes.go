package authentication

import (
	"github.com/gofiber/fiber/v2"
)

// SETUP ROUTES CONFIGURES AUTHENTICATION ENDPOINTS
func SetupRoutes(app fiber.Router, handler *Handler) {
	auth := app.Group("/authentications")

	auth.Post("/", handler.PostAuthenticationHandler)
	auth.Put("/", handler.PutAuthenticationHandler)
	auth.Delete("/", handler.DeleteAuthenticationHandler)
}
