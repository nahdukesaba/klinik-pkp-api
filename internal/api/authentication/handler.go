package authentication

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Handler struct {
	service *Service
}

// CONSTRUCTOR
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PostAuthenticationHandler(ctx *fiber.Ctx) error {
	var payload AddAuthenticationPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddRefreshToken(payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		if err == bcrypt.ErrMismatchedHashAndPassword {
			return utils.JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "authentication created", data, false)
}

// VALIDATES REFRESH TOKEN, GENERATES NEW ACCESS TOKEN
func (h *Handler) PutAuthenticationHandler(ctx *fiber.Ctx) error {
	var payload RefreshAuthenticationPayload

	// PARSE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.UpdateAccessToken(payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "access token refreshed", data, false)
}

// VALIDATES AND REMOVES REFRESH TOKEN FROM DATABASE
func (h *Handler) DeleteAuthenticationHandler(ctx *fiber.Ctx) error {
	var payload DeleteAuthenticationPayload

	// PARSE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.DeleteRefreshToken(payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "session terminated", nil, false)
}
