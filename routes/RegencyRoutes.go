package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRegencyRoutes configures all regency-related routes
func SetupRegencyRoutes(app fiber.Router, controller *controller.RegencyController) {
	// Public routes (read-only)
	regencies := app.Group("/regencies")
	regencies.Get("/", controller.GetAll)
	regencies.Get("/:id", controller.GetByID)

	// Admin routes (create, update, delete)
	regenciesAdmin := app.Group("/admin/regencies", middleware.RequireAuth(), middleware.RequireAdmin())
	regenciesAdmin.Post("/", controller.Create)
	regenciesAdmin.Put("/:id", controller.Update)
	regenciesAdmin.Delete("/:id", controller.Delete)
}
