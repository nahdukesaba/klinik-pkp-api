package routes

import (
	"klinik-pkp-api/controller"

	"github.com/gofiber/fiber/v2"
)

// SetupRegionRoutes configures region routes
func SetupRegionRoutes(app fiber.Router, controller *controller.RegionController) {
	app.Get("/regions", controller.GetAllRegions)
	app.Get("/region/:id", controller.GetRegionByID)
	app.Post("/region", controller.CreateRegion)
	app.Put("/region/:id", controller.UpdateRegion)
	app.Delete("/region/:id", controller.DeleteRegion)
}
