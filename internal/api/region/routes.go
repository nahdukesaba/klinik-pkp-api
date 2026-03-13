package region

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	regions := app.Group("/regions")

	regions.Get("/", handler.GetRegionsHandler)
	regions.Get("/:id", handler.GetRegionByIdHandler)
	regions.Post("/", handler.PostRegionHandler)
	regions.Put("/:id", handler.PutRegionByIdHandler)
	regions.Delete("/:id", handler.DeleteRegionByIdHandler)
}
