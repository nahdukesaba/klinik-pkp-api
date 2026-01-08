package region

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
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
	regions, err := h.service.GetRegions()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", regions))
}

func (h *Handler) GetRegionByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	region, err := h.service.GetRegionById(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", region))
}

func (h *Handler) PostRegionHandler(ctx *fiber.Ctx) error {
	var payload RegionPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	region, err := h.service.AddRegion(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("Region created", region))
}

func (h *Handler) PutRegionByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload RegionPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	region, err := h.service.EditRegionById(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Region updated", region))
}

func (h *Handler) DeleteRegionByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteRegionById(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Region deleted", nil))
}
