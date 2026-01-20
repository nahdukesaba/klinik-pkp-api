package province

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/provinces", handler.GetProvincesHandler)
	app.Get("/provinces/:id", handler.GetProvinceByIdHandler)
	app.Post("/provinces", handler.PostProvinceHandler)
	app.Put("/provinces/:id", handler.PutProvinceByIdHandler)
	app.Delete("/provinces/:id", handler.DeleteProvinceByIdHandler)
}
