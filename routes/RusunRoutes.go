package routes

import (
	"klinik-pkp-api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupRusunRoutes(app fiber.Router, controller *controller.RusunController) {
	app.Get("/rusuns", controller.GetAllRusun)
	app.Get("/rusun/:id", controller.GetRusunByID)
	app.Post("/rusun", controller.CreateRusun)
	app.Put("/rusun/:id", controller.UpdateRusun)
	app.Delete("/rusun/:id", controller.DeleteRusun)
}
