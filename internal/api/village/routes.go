package village

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/villages", handler.GetVillagesHandler)
	app.Get("/village/:id", handler.GetVillageByIdHandler)
	app.Post("/village", handler.PostVillageHandler)
	app.Put("/village/:id", handler.PutVillageByIdHandler)
	app.Delete("/village/:id", handler.DeleteVillageByIdHandler)
}
