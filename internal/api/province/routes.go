package province

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	provinces := app.Group("/provinces")

	provinces.Get("/", handler.GetProvincesHandler)
	provinces.Get("/:id", handler.GetProvinceByIdHandler)
	provinces.Post("/", handler.PostProvinceHandler)
	provinces.Put("/:id", handler.PutProvinceByIdHandler)
	provinces.Delete("/:id", handler.DeleteProvinceByIdHandler)
}
