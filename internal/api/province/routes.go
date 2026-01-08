package province

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/provinces", handler.GetProvincesHandler)
	app.Get("/province/:id", handler.GetProvinceByIdHandler)
	app.Post("/province", handler.PostProvinceHandler)
	app.Put("/province/:id", handler.PutProvinceByIdHandler)
	app.Delete("/province/:id", handler.DeleteProvinceByIdHandler)
}
