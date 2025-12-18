package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupBalaiRoutes configures all balai-related routes
func SetupBalaiRoutes(app fiber.Router, controller *controller.BalaiController) {
	// Public routes (read-only)
	balais := app.Group("/balais")
	balais.Get("/", controller.GetAll)
	balais.Get("/:id", controller.GetByID)

	// Admin routes (create, update, delete)
	balaisAdmin := app.Group("/admin/balais", middleware.RequireAuth(), middleware.RequireAdmin())
	balaisAdmin.Post("/", controller.Create)
	balaisAdmin.Put("/:id", controller.Update)
	balaisAdmin.Delete("/:id", controller.Delete)
}
