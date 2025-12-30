package controller

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"
)

// REGION CONTROLLER STRUCT
type RegionController struct {
	service *service.RegionService
}

// REGION CONTROLLER CONSTRUCTOR
func NewRegionController(service *service.RegionService) *RegionController {
	return &RegionController{service: service}
}

func (c *RegionController) GetAllRegions(ctx *fiber.Ctx) error {
	regions, err := c.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch regions"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", regions))
}

func (c *RegionController) GetRegionByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	region, err := c.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", region))
}

func (c *RegionController) CreateRegion(ctx *fiber.Ctx) error {
	var input service.RegionPayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	region, err := c.service.Create(input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Region created", region))
}

func (c *RegionController) UpdateRegion(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var input service.RegionPayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	region, err := c.service.Update(id, input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Region updated", region))
}

func (c *RegionController) DeleteRegion(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Region deleted", nil))
}
