package kumuh

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	kumuh := app.Group("/kumuh")
	
	kumuh.Get("/", handler.GetKumuhHandler)
	kumuh.Get("/:id", handler.GetKumuhByIdHandler)
	kumuh.Post("/", handler.PostKumuhHandler)
	kumuh.Put("/:id", handler.PutKumuhByIdHandler)
	kumuh.Delete("/:id", handler.DeleteKumuhByIdHandler)
}
