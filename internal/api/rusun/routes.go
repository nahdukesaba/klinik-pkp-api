package rusun

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/rusun", handler.GetRusunHandler)
	app.Get("/rusun/:id", handler.GetRusunByIdHandler)
	app.Post("/rusun", handler.PostRusunHandler)
	app.Put("/rusun/:id", handler.PutRusunByIdHandler)
	app.Delete("/rusun/:id", handler.DeleteRusunByIdHandler)
}
