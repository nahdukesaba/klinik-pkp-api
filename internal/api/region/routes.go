package region

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures region routes
func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/regions", handler.GetAllRegions)
	app.Get("/region/:id", handler.GetRegionByID)
	app.Post("/region", handler.CreateRegion)
	app.Put("/region/:id", handler.UpdateRegion)
	app.Delete("/region/:id", handler.DeleteRegion)
}
