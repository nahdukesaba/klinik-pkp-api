package user

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures user profile routes (protected)
func SetupRoutes(app fiber.Router, handler *Handler) {
	// Users list endpoint
	app.Get("/users", handler.GetAllUsers)
	// app.Get("/province/:id", handler.GetUserByID)
	// app.Post("/province", handler.CreateUser)
	// app.Put("/province/:id", handler.UpdateUser)
	// app.Delete("/province/:id", handler.DeleteUser)
}
