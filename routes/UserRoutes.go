package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures user profile routes (protected)
func SetupUserRoutes(app fiber.Router, controller *controller.UserController) {
	// Protected routes (authentication required)
	user := app.Group("/user", middleware.RequireAuth())
	user.Get("/profile", controller.GetProfile)
	user.Put("/profile", controller.UpdateProfile)
	user.Put("/change-password", controller.ChangePassword)
}
