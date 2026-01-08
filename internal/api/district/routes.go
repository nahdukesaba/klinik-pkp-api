package district

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/districts", handler.GetDistrictsHandler)
	app.Get("/district/:id", handler.GetDistrictByIDHandler)
	app.Post("/district", handler.PostDistrictHandler)
	app.Put("/district/:id", handler.PutDistrictByIdHandler)
	app.Delete("/district/:id", handler.DeleteDistrictByIdHandler)
}
