package controller

import (
	"klinik-api/service"
	"klinik-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	service *service.AdminService
}

func NewAdminController(service *service.AdminService) *AdminController {
	return &AdminController{service: service}
}

// Login handles admin authentication
func (c *AdminController) Login(ctx *fiber.Ctx) error {
	var input service.AdminLoginInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body: " + err.Error()))
	}

	result, err := c.service.Login(input)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Login successful", result))
}

// GetAllUsers returns all users with pagination
func (c *AdminController) GetAllUsers(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := c.service.GetAllUsers(page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.Response{
		Success: true,
		Message: "Success",
		Data: fiber.Map{
			"users": users,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetUserByID returns user by ID
func (c *AdminController) GetUserByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid user ID"))
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", user))
}

// CreateUser creates new user (admin only)
func (c *AdminController) CreateUser(ctx *fiber.Ctx) error {
	var input service.CreateUserInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	// Get admin user ID from context
	adminID := ctx.Locals("userID").(uint)

	user, err := c.service.CreateUser(input, adminID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("User created successfully", user))
}

// UpdateUser updates user (admin only)
func (c *AdminController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid user ID"))
	}

	var input service.UpdateUserInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	// Get admin user ID from context
	adminID := ctx.Locals("userID").(uint)

	user, err := c.service.UpdateUser(uint(id), input, adminID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("User updated successfully", user))
}

// DeleteUser deletes user (admin only)
func (c *AdminController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil || id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid user ID"))
	}

	// Get admin user ID from context
	adminID := ctx.Locals("userID").(uint)

	if err := c.service.DeleteUser(uint(id), adminID); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("User deleted successfully", nil))
}
