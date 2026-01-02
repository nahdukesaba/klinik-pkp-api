package village

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures village routes
func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/villages", handler.GetAllVillages)
	app.Get("/village/:id", handler.GetVillageByID)
	app.Post("/village", handler.CreateVillage)
	app.Put("/village/:id", handler.UpdateVillage)
	app.Delete("/village/:id", handler.DeleteVillage)
}
