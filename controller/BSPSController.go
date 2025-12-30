package controller

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"
)

// BSPS CONTROLLER STRUCT
type BSPSController struct {
	service *service.BSPSService
}

// BSPS CONTROLLER CONSTRUCTOR
func NewBSPSController(service *service.BSPSService) *BSPSController {
	return &BSPSController{service: service}
}

func (c *BSPSController) GetAllBSPS(ctx *fiber.Ctx) error {
	bsps, err := c.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch bsps data"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", bsps))
}
