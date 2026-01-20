package region

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/regions", handler.GetRegionsHandler)
	app.Get("/regions/:id", handler.GetRegionByIdHandler)
	app.Post("/regions", handler.PostRegionHandler)
	app.Put("/regions/:id", handler.PutRegionByIdHandler)
	app.Delete("/regions/:id", handler.DeleteRegionByIdHandler)
}
