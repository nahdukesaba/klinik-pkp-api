package kumuh

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	kumuh := app.Group("/kumuh")

	kumuh.Get("/", handler.GetKumuhHandler)
	kumuh.Get("/years", handler.GetAvailableYearsHandler)
	kumuh.Get("/:id", handler.GetKumuhByIdHandler)
	kumuh.Post("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostKumuhHandler)
	kumuh.Put("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutKumuhByIdHandler)
	kumuh.Delete("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteKumuhByIdHandler)
}
