package district

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/districts", handler.GetDistrictsHandler)
	app.Get("/districts/:id", handler.GetDistrictByIDHandler)
	app.Post("/districts", handler.PostDistrictHandler)
	app.Put("/districts/:id", handler.PutDistrictByIdHandler)
	app.Delete("/districts/:id", handler.DeleteDistrictByIdHandler)
}
