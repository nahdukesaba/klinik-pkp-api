package routes

import (
	"klinik-pkp-api/controller"

	"github.com/gofiber/fiber/v2"
)

// SetupProvinceRoutes configures province routes
func SetupProvinceRoutes(app fiber.Router, controller *controller.ProvinceController) {
	app.Get("/provinces", controller.GetAllProvinces)
	app.Get("/province/:id", controller.GetProvinceByID)
	app.Post("/province", controller.CreateProvince)
	app.Put("/province/:id", controller.UpdateProvince)
	app.Delete("/province/:id", controller.DeleteProvince)
}
