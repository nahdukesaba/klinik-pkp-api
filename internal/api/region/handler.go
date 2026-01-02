package region

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler struct
type Handler struct {
	service *Service
}

// NewHandler constructor
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllRegions(ctx *fiber.Ctx) error {
	regions, err := h.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch regions"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", regions))
}

func (h *Handler) GetRegionByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	region, err := h.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", region))
}

func (h *Handler) CreateRegion(ctx *fiber.Ctx) error {
	var payload RegionPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	region, err := h.service.Create(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Region created", region))
}

func (h *Handler) UpdateRegion(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload RegionPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	region, err := h.service.Update(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Region updated", region))
}

func (h *Handler) DeleteRegion(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Region deleted", nil))
}
