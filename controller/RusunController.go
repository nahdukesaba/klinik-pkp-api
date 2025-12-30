package controller

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"
	"strconv"
)

// RUSUN CONTROLLER STRUCT
type RusunController struct {
	service *service.RusunService
}

// RUSUN CONTROLLER CONSTRUCTOR
func NewRusunController(service *service.RusunService) *RusunController {
	return &RusunController{service: service}
}

func (c *RusunController) GetAllRusun(ctx *fiber.Ctx) error {
	rusun, err := c.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch rusun data"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", rusun))
}

func (c *RusunController) GetRusunByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid rusun id"))
	}

	rusun, err := c.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", rusun))
}

func (c *RusunController) CreateRusun(ctx *fiber.Ctx) error {
	var input service.RusunPayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	rusun, err := c.service.Create(input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Rusun created", rusun))
}

func (c *RusunController) UpdateRusun(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid rusun id"))
	}

	var input service.RusunPayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	rusun, err := c.service.Update(id, input)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Rusun updated", rusun))
}

func (c *RusunController) DeleteRusun(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid rusun id"))
	}

	err = c.service.Delete(id)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Rusun deleted", nil))
}
