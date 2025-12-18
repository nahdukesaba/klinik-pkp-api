package controller

import (
	"klinik-api/service"
	"klinik-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BalaiController struct {
	service *service.BalaiService
}

func NewBalaiController(service *service.BalaiService) *BalaiController {
	return &BalaiController{service: service}
}

// GetAll returns all balais
func (c *BalaiController) GetAll(ctx *fiber.Ctx) error {
	provinceID, _ := strconv.ParseUint(ctx.Query("province_id"), 10, 32)
	categoryID, _ := strconv.ParseUint(ctx.Query("category_id"), 10, 32)

	balais, err := c.service.GetAll(uint(provinceID), uint(categoryID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", balais))
}

// GetByID returns balai by ID
func (c *BalaiController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid balai ID"))
	}

	balai, err := c.service.GetByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", balai))
}

// Create creates new balai (admin only)
func (c *BalaiController) Create(ctx *fiber.Ctx) error {
	var input struct {
		ProvinceID uint   `json:"province_id"`
		CategoryID uint   `json:"category_id"`
		Name       string `json:"name"`
		Address    string `json:"address"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	balai, err := c.service.Create(input.ProvinceID, input.CategoryID, input.Name, input.Address)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Balai created successfully", balai))
}

// Update updates balai (admin only)
func (c *BalaiController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid balai ID"))
	}

	var input struct {
		ProvinceID uint   `json:"province_id"`
		CategoryID uint   `json:"category_id"`
		Name       string `json:"name"`
		Address    string `json:"address"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	balai, err := c.service.Update(uint(id), input.ProvinceID, input.CategoryID, input.Name, input.Address)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Balai updated successfully", balai))
}

// Delete deletes balai (admin only)
func (c *BalaiController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid balai ID"))
	}

	if err := c.service.Delete(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Balai deleted successfully", nil))
}
