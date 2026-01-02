package province

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures province routes
func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/provinces", handler.GetAllProvinces)
	app.Get("/province/:id", handler.GetProvinceByID)
	app.Post("/province", handler.CreateProvince)
	app.Put("/province/:id", handler.UpdateProvince)
	app.Delete("/province/:id", handler.DeleteProvince)
}
