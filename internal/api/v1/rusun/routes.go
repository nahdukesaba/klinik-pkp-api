package rusun

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	rusun := app.Group("/rusun")

	rusun.Get("/", handler.GetRusunHandler)
	rusun.Get("/years", handler.GetAvailableYearsHandler)
	rusun.Get("/:id", handler.GetRusunByIdHandler)
	rusun.Post("/", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PostRusunHandler)
	rusun.Put("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PutRusunByIdHandler)
	rusun.Delete("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.DeleteRusunByIdHandler)
}
