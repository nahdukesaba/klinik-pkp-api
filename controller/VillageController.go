package controller

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"
)

// VILLAGE CONTROLLER STRUCT
type VillageController struct {
	service *service.VillageService
}

// VILLAGE CONTROLLER CONSTRUCTOR
func NewVillageController(service *service.VillageService) *VillageController {
	return &VillageController{service: service}
}

func (c *VillageController) GetAllVillages(ctx *fiber.Ctx) error {
	villages, err := c.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch villages"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", villages))
}

func (c *VillageController) GetVillageByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	village, err := c.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", village))
}

func (c *VillageController) CreateVillage(ctx *fiber.Ctx) error {
	var input service.VillagePayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	village, err := c.service.Create(input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Village created", village))
}

func (c *VillageController) UpdateVillage(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var input service.VillagePayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	village, err := c.service.Update(id, input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Village updated", village))
}

func (c *VillageController) DeleteVillage(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Village deleted", nil))
}
