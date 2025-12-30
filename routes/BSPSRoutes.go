package routes

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/controller"
)

func SetupBSPSRoutes(app fiber.Router, controller *controller.BSPSController) {
	app.Get("/bsps", controller.GetAllBSPS)
}
