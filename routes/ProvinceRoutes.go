package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupProvinceRoutes configures all province-related routes
func SetupProvinceRoutes(app fiber.Router, controller *controller.ProvinceController) {
	// Public routes (read-only)
	provinces := app.Group("/provinces")
	provinces.Get("/", controller.GetAll)
	provinces.Get("/:id", controller.GetByID)

	// Admin routes (create, update, delete)
	provincesAdmin := app.Group("/admin/provinces", middleware.RequireAuth(), middleware.RequireAdmin())
	provincesAdmin.Post("/", controller.Create)
	provincesAdmin.Put("/:id", controller.Update)
	provincesAdmin.Delete("/:id", controller.Delete)
}
