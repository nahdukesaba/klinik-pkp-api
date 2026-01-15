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
	bsps, err := h.service.GetBSPS()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", bsps))
}

func (h *Handler) GetBSPSByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	bsps, err := h.service.GetBSPSById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", bsps))
}

func (h *Handler) PostBSPSHandler(ctx *fiber.Ctx) error {
	var payload BSPSPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	bsps, err := h.service.AddBSPS(payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("BSPS created", bsps))
}

func (h *Handler) PutBSPSByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	var payload BSPSPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	bsps, err := h.service.EditBSPSById(id, payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("BSPS updated", bsps))
}

func (h *Handler) DeleteBSPSByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	err = h.service.DeleteBSPSById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("BSPS deleted", nil))
}
