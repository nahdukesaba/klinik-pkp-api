package province

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

func (h *Handler) GetAllProvinces(ctx *fiber.Ctx) error {
	provinces, err := h.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("success", provinces))
}

func (h *Handler) GetProvinceByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	province, err := h.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", province))
}

func (h *Handler) CreateProvince(ctx *fiber.Ctx) error {
	var payload ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	province, err := h.service.Create(payload)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}
	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("Province created", province))
}

func (h *Handler) UpdateProvince(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	province, err := h.service.Update(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Province updated", province))
}

func (h *Handler) DeleteProvince(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Province deleted", nil))
}
