package controller

import (
	"klinik-pkp-api/service"
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

// Register handles user registration
func (c *UserController) Register(ctx *fiber.Ctx) error {
	var input service.RegisterPayload
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body: " + err.Error()))
	}

	result, err := c.service.Register(input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Registration successful", result))
}

// Login handles user authentication
func (c *UserController) Login(ctx *fiber.Ctx) error {
	var input service.LoginPayload
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body: " + err.Error()))
	}

	result, err := c.service.Login(input)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Login successful", result))
}

// GetAllUsers returns all users in the database
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.service.GetAll()
	
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch users"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", users))
}

// GetProfile returns current user profile
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	user, err := c.service.GetProfile(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("User not found"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", user))
}

// UpdateProfile updates current user profile
func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var input struct {
		Name string `json:"name"`
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	user, err := c.service.UpdateProfile(userID, input.Name)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Profile updated successfully", user))
}

// ChangePassword handles password change
func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var input service.ChangePasswordInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := c.service.ChangePassword(userID, input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Password changed successfully", nil))
}
