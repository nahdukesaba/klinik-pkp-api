package province

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	provinces := app.Group("/provinces")

	provinces.Get("/", handler.GetProvincesHandler)
	provinces.Get("/:id", handler.GetProvinceByIdHandler)
	provinces.Post("/", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PostProvinceHandler)
	provinces.Put("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PutProvinceByIdHandler)
	provinces.Delete("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.DeleteProvinceByIdHandler)
}
