package routes

import (
	"klinik-pkp-api/controller"

	"github.com/gofiber/fiber/v2"
)

// SetupVillageRoutes configures village routes
func SetupVillageRoutes(app fiber.Router, controller *controller.VillageController) {
	app.Get("/villages", controller.GetAllVillages)
	app.Get("/village/:id", controller.GetVillageByID)
	app.Post("/village", controller.CreateVillage)
	app.Put("/village/:id", controller.UpdateVillage)
	app.Delete("/village/:id", controller.DeleteVillage)
}
