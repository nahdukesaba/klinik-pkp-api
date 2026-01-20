package village

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/villages", handler.GetVillagesHandler)
	app.Get("/villages/:id", handler.GetVillageByIdHandler)
	app.Post("/villages", handler.PostVillageHandler)
	app.Put("/villages/:id", handler.PutVillageByIdHandler)
	app.Delete("/villages/:id", handler.DeleteVillageByIdHandler)
}
