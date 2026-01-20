package bsps

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
	"strconv"
)

// CURRENT INSTANCE
type Handler struct {
	service *Service
}

// CONSTRUCTOR
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetBSPSHandler(ctx *fiber.Ctx) error {
	data, err := h.service.GetBSPS()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) GetBSPSByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	data, err := h.service.GetBSPSById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) PostBSPSHandler(ctx *fiber.Ctx) error {
	var payload BSPSPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	data, err := h.service.AddBSPS(payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "BSPS created", data, false)
}

func (h *Handler) PutBSPSByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	var payload BSPSPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	data, err := h.service.EditBSPSById(id, payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "BSPS updated", data, false)
}

func (h *Handler) DeleteBSPSByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	err = h.service.DeleteBSPSById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "BSPS deleted", nil, false)
}
