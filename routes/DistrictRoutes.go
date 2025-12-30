package routes

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/controller"
)

// SetupDistrictRoutes configures district routes
func SetupDistrictRoutes(app fiber.Router, controller *controller.DistrictController) {
	app.Get("/districts", controller.GetAllDistricts)
	app.Get("/district/:id", controller.GetDistrictByID)
	app.Post("/district", controller.CreateDistrict)
	app.Put("/district/:id", controller.UpdateDistrict)
	app.Delete("/district/:id", controller.DeleteDistrict)
}
