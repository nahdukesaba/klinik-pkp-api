package region

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/regions", handler.GetRegionsHandler)
	app.Get("/region/:id", handler.GetRegionByIdHandler)
	app.Post("/region", handler.PostRegionHandler)
	app.Put("/region/:id", handler.PutRegionByIdHandler)
	app.Delete("/region/:id", handler.DeleteRegionByIdHandler)
}
