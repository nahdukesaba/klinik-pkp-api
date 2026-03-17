package village

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	villages := app.Group("/villages")

	villages.Get("/", handler.GetVillagesHandler)
	villages.Get("/:id", handler.GetVillageByIdHandler)
	villages.Post("/", handler.PostVillageHandler)
	villages.Put("/:id", handler.PutVillageByIdHandler)
	villages.Delete("/:id", handler.DeleteVillageByIdHandler)
}
