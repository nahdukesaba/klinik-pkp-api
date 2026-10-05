package village

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	villages := app.Group("/villages")

	villages.Get("/", handler.GetVillagesHandler)
	villages.Get("/:id", handler.GetVillageByIdHandler)
	villages.Post("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostVillageHandler)
	villages.Put("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutVillageByIdHandler)
	villages.Delete("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteVillageByIdHandler)
}
