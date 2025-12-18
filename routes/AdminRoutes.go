package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes configures admin-related routes (protected)
func SetupAdminRoutes(app fiber.Router, controller *controller.AdminController) {
	// Protected admin routes (authentication + admin role required)
	admin := app.Group("/admin", middleware.RequireAuth(), middleware.RequireAdmin())
	
	// User management
	admin.Get("/users", controller.GetAllUsers)
	admin.Get("/users/:id", controller.GetUserByID)
	admin.Post("/users", controller.CreateUser)
	admin.Put("/users/:id", controller.UpdateUser)
	admin.Delete("/users/:id", controller.DeleteUser)
}
