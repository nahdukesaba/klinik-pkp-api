package routes

import (
	"klinik-api/controller"
	"klinik-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRusunRoutes sets up all rusun routes
func SetupRusunRoutes(app fiber.Router, ctrl *controller.RusunController) {
	rusun := app.Group("/rusun")
	{
		rusun.Get("/", ctrl.GetAll)
		rusun.Get("/:id", ctrl.GetByID)
		rusun.Get("/province/:province_id", ctrl.GetByProvince)
		rusun.Get("/map", ctrl.GetMapData)
		rusun.Get("/statistics", ctrl.GetStatistics)
	}

	admin := app.Group("/admin/rusun", middleware.RequireAuth(), middleware.RequireAdmin())
	{
		admin.Post("/", ctrl.Create)
		admin.Put("/:id", ctrl.Update)
		admin.Delete("/:id", ctrl.Delete)
	}
}
