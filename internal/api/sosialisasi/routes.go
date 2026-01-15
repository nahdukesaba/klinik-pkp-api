package sosialisasi

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/sosialisasi", handler.GetSosialisasiHandler)
	app.Get("/sosialisasi/:id", handler.GetSosialisasiByIdHandler)
	app.Post("/sosialisasi", handler.AddSosialisasiHandler)
	app.Put("/sosialisasi/:id", handler.PutSosialisasiByIdHandler)
	app.Delete("/sosialisasi/:id", handler.DeleteSosialisasiByIdHandler)
}