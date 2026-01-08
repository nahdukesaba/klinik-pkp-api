package sosialisasi

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/sosialisasi", handler.GetSosialisasiHandler)
}