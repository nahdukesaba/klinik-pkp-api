package controller

import (
	"klinik-api/service"
	"klinik-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RegencyController struct {
	service *service.RegencyService
}

func NewRegencyController(service *service.RegencyService) *RegencyController {
	return &RegencyController{service: service}
}

// GetAll returns all regencies
func (c *RegencyController) GetAll(ctx *fiber.Ctx) error {
	provinceID, _ := strconv.ParseUint(ctx.Query("province_id"), 10, 32)
	includeProvince := ctx.Query("include_province") == "true"

	regencies, err := c.service.GetAll(uint(provinceID), includeProvince)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", regencies))
}

// GetByID returns regency by ID
func (c *RegencyController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid regency ID"))
	}

	includeProvince := ctx.Query("include_province") == "true"

	regency, err := c.service.GetByID(uint(id), includeProvince)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", regency))
}

// Create creates new regency (admin only)
func (c *RegencyController) Create(ctx *fiber.Ctx) error {
	var input struct {
		ProvinceID uint   `json:"province_id"`
		Name       string `json:"name"`
		Code       string `json:"code"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	regency, err := c.service.Create(input.ProvinceID, input.Name, input.Code)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Regency created successfully", regency))
}

// Update updates regency (admin only)
func (c *RegencyController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid regency ID"))
	}

	var input struct {
		ProvinceID uint   `json:"province_id"`
		Name       string `json:"name"`
		Code       string `json:"code"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	regency, err := c.service.Update(uint(id), input.ProvinceID, input.Name, input.Code)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Regency updated successfully", regency))
}

// Delete deletes regency (admin only)
func (c *RegencyController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid regency ID"))
	}

	if err := c.service.Delete(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Regency deleted successfully", nil))
}
