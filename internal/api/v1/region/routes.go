package region

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	regions := app.Group("/regions")

	regions.Get("/", handler.GetRegionsHandler)
	regions.Get("/:id", handler.GetRegionByIdHandler)
	regions.Post("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostRegionHandler)
	regions.Put("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutRegionByIdHandler)
	regions.Delete("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteRegionByIdHandler)
}
