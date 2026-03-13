package user

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	users := app.Group("/users")

	// PUBLIC ROUTE - REGISTRATION ONLY
	// users.Post("/register", handler.Register)

	// PROTECTED ROUTES (REQUIRES AUTH MIDDLEWARE)
	users.Get("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.GetUsersHandler)
	users.Get("/:id", utils.AuthMiddleware(),  handler.GetUserByIdHandler)
	users.Put("/:id", utils.AuthMiddleware(), handler.PutUserByIdHandler)
}
