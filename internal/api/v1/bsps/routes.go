package bsps

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	bsps := app.Group("/bsps")

	bsps.Get("/", handler.GetBSPSHandler)
	bsps.Get("/years", handler.GetAvailableYearsHandler)
	bsps.Get("/:id", handler.GetBSPSByIdHandler)
	bsps.Post("/", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PostBSPSHandler)
	bsps.Put("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PutBSPSByIdHandler)
	bsps.Delete("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.DeleteBSPSByIdHandler)
}
