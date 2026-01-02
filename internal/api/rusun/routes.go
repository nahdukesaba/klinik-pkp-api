package rusun

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/rusun", handler.GetAllRusun)
	app.Get("/rusun/:id", handler.GetRusunByID)
	app.Post("/rusun", handler.CreateRusun)
	app.Put("/rusun/:id", handler.UpdateRusun)
	app.Delete("/rusun/:id", handler.DeleteRusun)
}
