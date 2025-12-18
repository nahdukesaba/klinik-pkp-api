package routes

import (
	"klinik-api/controller"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures all authentication routes (public endpoints)
func SetupAuthRoutes(app fiber.Router, userCtrl *controller.UserController, adminCtrl *controller.AdminController) {
	// Public user authentication endpoints (no auth required)
	app.Post("/register", userCtrl.Register) // User self-registration with default "User" role
	app.Post("/login", userCtrl.Login)       // User login returns JWT token

	// Public admin authentication endpoint (no auth required)
	admin := app.Group("/admin")
	admin.Post("/login", adminCtrl.Login) // Admin login returns JWT token
}
