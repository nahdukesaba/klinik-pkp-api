package kumuh

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/kumuh", handler.GetKumuhHandler)
	app.Get("/kumuh/:id", handler.GetKumuhByIdHandler)
	app.Post("/kumuh", handler.PostKumuhHandler)
	app.Put("/kumuh/:id", handler.PutKumuhByIdHandler)
	app.Delete("/kumuh/:id", handler.DeleteKumuhByIdHandler)
}
