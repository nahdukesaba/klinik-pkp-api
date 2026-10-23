package user

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	users := app.Group("/users")

	// PUBLIC ROUTE - REGISTRATION ONLY
	// users.Post("/register", handler.Register)

	// PROTECTED ROUTES (REQUIRES AUTH MIDDLEWARE)
	users.Get("/", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.GetUsersHandler)
	users.Get("/me", utils.AuthMiddleware(handler.service.db), handler.GetCurrentUserHandler)
	users.Get("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.GetUserByIdHandler)
	users.Post("/", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PostUserHandler)
	users.Put("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.PutUserByIdHandler)
	users.Delete("/:id", utils.AuthMiddleware(handler.service.db), utils.RoleMiddleware("admin"), handler.DeleteUserByIdHandler)
}
