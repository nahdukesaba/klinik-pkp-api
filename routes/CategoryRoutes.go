package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupCategoryRoutes configures all category-related routes
func SetupCategoryRoutes(app fiber.Router, controller *controller.CategoryController) {
	// Public routes (read-only)
	categories := app.Group("/categories")
	categories.Get("/", controller.GetAll)
	categories.Get("/:id", controller.GetByID)

	// Admin routes (create, update, delete)
	categoriesAdmin := app.Group("/admin/categories", middleware.RequireAuth(), middleware.RequireAdmin())
	categoriesAdmin.Post("/", controller.Create)
	categoriesAdmin.Put("/:id", controller.Update)
	categoriesAdmin.Delete("/:id", controller.Delete)
}
