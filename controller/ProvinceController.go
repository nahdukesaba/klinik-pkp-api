package controller

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"
)

// PROVINCE CONTROLLER STRUCT
type ProvinceController struct {
	service *service.ProvinceService
}

// PROVINCE CONTROLLER CONSTRUCTOR
func NewProvinceController(service *service.ProvinceService) *ProvinceController {
	return &ProvinceController{service: service}
}

func (c *ProvinceController) GetAllProvinces(ctx *fiber.Ctx) error {
	provinces, err := c.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch provinces"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("success", provinces))
}

func (c *ProvinceController) GetProvinceByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	province, err := c.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", province))
}

func (c *ProvinceController) CreateProvince(ctx *fiber.Ctx) error {
	var input service.ProvincePayload

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	province, err := c.service.Create(input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Province created", province))
}

func (c *ProvinceController) UpdateProvince(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload service.ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	province, err := c.service.Update(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Province updated", province))
}

func (c *ProvinceController) DeleteProvince(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Province deleted", nil))
}
