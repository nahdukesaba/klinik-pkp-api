package sosialisasi

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	sosialisasi := app.Group("/sosialisasi")
	
	sosialisasi.Get("/", handler.GetSosialisasiHandler)
	sosialisasi.Get("/:id", handler.GetSosialisasiByIdHandler)
	sosialisasi.Post("/", handler.PostSosialisasiHandler)
	sosialisasi.Put("/:id", handler.PutSosialisasiByIdHandler)
	sosialisasi.Delete("/:id", handler.DeleteSosialisasiByIdHandler)
}