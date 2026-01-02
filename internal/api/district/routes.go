package district

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures district routes
func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/districts", handler.GetAllDistricts)
	app.Get("/district/:id", handler.GetDistrictByID)
	app.Post("/district", handler.CreateDistrict)
	app.Put("/district/:id", handler.UpdateDistrict)
	app.Delete("/district/:id", handler.DeleteDistrict)
}
