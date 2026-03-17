package district

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	districts := app.Group("/districts")
	
	districts.Get("/", handler.GetDistrictsHandler)
	districts.Get("/:id", handler.GetDistrictByIDHandler)
	districts.Post("/", handler.PostDistrictHandler)
	districts.Put("/:id", handler.PutDistrictByIdHandler)
	districts.Delete("/:id", handler.DeleteDistrictByIdHandler)
}
