package rusun

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	rusun := app.Group("/rusun")

	rusun.Get("/", handler.GetRusunHandler)
	rusun.Get("/:id", handler.GetRusunByIdHandler)
	rusun.Post("/", handler.PostRusunHandler)
	rusun.Put("/:id", handler.PutRusunByIdHandler)
	rusun.Delete("/:id", handler.DeleteRusunByIdHandler)
}
