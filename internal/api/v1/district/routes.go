package district

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	districts := app.Group("/districts")

	districts.Get("/", handler.GetDistrictsHandler)
	districts.Get("/:id", handler.GetDistrictByIDHandler)
	districts.Post("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostDistrictHandler)
	districts.Put("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutDistrictByIdHandler)
	districts.Delete("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteDistrictByIdHandler)
}
