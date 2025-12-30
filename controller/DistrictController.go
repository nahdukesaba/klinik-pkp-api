package controller

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"
)

// DISTRICT CONTROLLER STRUCT
type DistrictController struct {
	service *service.DistrictService
}

// DISTRICT CONTROLLER CONSTRUCTOR
func NewDistrictController(service *service.DistrictService) *DistrictController {
	return &DistrictController{service: service}
}

func (c *DistrictController) GetAllDistricts(ctx *fiber.Ctx) error {
	districts, err := c.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch districts"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("success", districts))
}

func (c *DistrictController) GetDistrictByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	district, err := c.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", district))
}

func (c *DistrictController) CreateDistrict(ctx *fiber.Ctx) error {
	var input service.DistrictPayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	district, err := c.service.Create(input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("District created", district))
}

func (c *DistrictController) UpdateDistrict(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var input service.DistrictPayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	district, err := c.service.Update(id, input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("District updated", district))
}

func (c *DistrictController) DeleteDistrict(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("District deleted", nil))
}
