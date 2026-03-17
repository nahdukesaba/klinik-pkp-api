package uploads

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(api fiber.Router) {
	// GROUP ROUTES WITH /api/v1/uploads PREFIX
	uploads := api.Group("/uploads")

	// SERVE STATIC FILES
	uploads.Static("/", "./storage")
}
