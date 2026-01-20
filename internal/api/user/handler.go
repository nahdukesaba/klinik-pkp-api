package user

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles user registration
func (h *Handler) Register(ctx *fiber.Ctx) error {
	var payload RegisterPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	result, err := h.service.Register(payload)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "Registration successful", result, false)
}

// Login handles user authentication
func (h *Handler) Login(ctx *fiber.Ctx) error {
	var payload LoginPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	result, err := h.service.Login(payload)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusUnauthorized, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "Login successful", result, false)
}

// GetAllUsers returns all users in the database
func (h *Handler) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := h.service.GetAll()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "Success", users, false)
}

// GetProfile returns current user profile
func (h *Handler) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	user, err := h.service.GetProfile(userID)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "Success", user, false)
}

// UpdateProfile updates current user profile
func (h *Handler) UpdateProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var payload struct {
		Name string `json:"name"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	user, err := h.service.UpdateProfile(userID, payload.Name)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "Profile updated successfully", user, false)
}
