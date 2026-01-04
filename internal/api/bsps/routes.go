package bsps

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/bsps", handler.GetBSPSHandler)
	app.Get("/bsps/:id", handler.GetBSPSByIdHandler)
	app.Post("/bsps", handler.PostBSPSHandler)
	app.Put("/bsps/:id", handler.PutBSPSByIdHandler)
	app.Delete("/bsps/:id", handler.DeleteBSPSByIdHandler)
}
