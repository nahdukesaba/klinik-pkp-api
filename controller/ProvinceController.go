package controller

import (
	"klinik-api/service"
	"klinik-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProvinceController struct {
	service *service.ProvinceService
}

func NewProvinceController(service *service.ProvinceService) *ProvinceController {
	return &ProvinceController{service: service}
}

// GetAll returns all provinces
func (c *ProvinceController) GetAll(ctx *fiber.Ctx) error {
	includeRegencies := ctx.Query("include_regencies") == "true"

	provinces, err := c.service.GetAll(includeRegencies)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", provinces))
}

// GetByID returns province by ID
func (c *ProvinceController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid province ID"))
	}

	includeRegencies := ctx.Query("include_regencies") == "true"

	province, err := c.service.GetByID(uint(id), includeRegencies)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", province))
}

// Create creates new province (admin only)
func (c *ProvinceController) Create(ctx *fiber.Ctx) error {
	var input struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	province, err := c.service.Create(input.Name, input.Code)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Province created successfully", province))
}

// Update updates province (admin only)
func (c *ProvinceController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid province ID"))
	}

	var input struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	province, err := c.service.Update(uint(id), input.Name, input.Code)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Province updated successfully", province))
}

// Delete deletes province (admin only)
func (c *ProvinceController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid province ID"))
	}

	if err := c.service.Delete(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Province deleted successfully", nil))
}
