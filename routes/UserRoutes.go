package routes

import (
	"klinik-pkp-api/controller"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures user profile routes (protected)
func SetupUserRoutes(app fiber.Router, controller *controller.UserController) {
	// Users list endpoint
	app.Get("/users", controller.GetAllUsers)
	// app.Get("/province/:id", controller.GetUserByID)
	// app.Post("/province", controller.CreateUser)
	// app.Put("/province/:id", controller.UpdateUser)
	// app.Delete("/province/:id", controller.DeleteUser)
}
