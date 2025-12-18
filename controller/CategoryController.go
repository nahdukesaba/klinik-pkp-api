package controller

import (
	"klinik-api/service"
	"klinik-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	service *service.CategoryService
}

func NewCategoryController(service *service.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// GetAll returns all categories
func (c *CategoryController) GetAll(ctx *fiber.Ctx) error {
	categories, err := c.service.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", categories))
}

// GetByID returns category by ID
func (c *CategoryController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid category ID"))
	}

	category, err := c.service.GetByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", category))
}

// Create creates new category (admin only)
func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	category, err := c.service.Create(input.Name, input.Description)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Category created successfully", category))
}

// Update updates category (admin only)
func (c *CategoryController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid category ID"))
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	category, err := c.service.Update(uint(id), input.Name, input.Description)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Category updated successfully", category))
}

// Delete deletes category (admin only)
func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid category ID"))
	}

	if err := c.service.Delete(uint(id)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Category deleted successfully", nil))
}
