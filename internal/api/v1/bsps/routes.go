package bsps

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	bsps := app.Group("/bsps")

	bsps.Get("/", handler.GetBSPSHandler)
	bsps.Get("/:id", handler.GetBSPSByIdHandler)
	bsps.Post("/", handler.PostBSPSHandler)
	bsps.Put("/:id", handler.PutBSPSByIdHandler)
	bsps.Delete("/:id", handler.DeleteBSPSByIdHandler)
}
