package region

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

// CURRENT INSTANCE
type Handler struct {
	service *Service
}

// CONSTRUCTOR
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRegionsHandler(ctx *fiber.Ctx) error {
	response, err := h.service.GetRegions()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", response))
}

func (h *Handler) GetRegionByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	response, err := h.service.GetRegionById(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", response))
}

func (h *Handler) PostRegionHandler(ctx *fiber.Ctx) error {
	var payload RegionPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	response, err := h.service.AddRegion(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("region created", response))
}

func (h *Handler) PutRegionByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload RegionPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	response, err := h.service.EditRegionById(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Region updated", response))
}

func (h *Handler) DeleteRegionByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteRegionById(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("region deleted", nil))
}
