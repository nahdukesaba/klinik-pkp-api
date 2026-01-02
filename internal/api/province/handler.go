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
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch provinces"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("success", provinces))
}

func (h *Handler) GetProvinceByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	province, err := h.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", province))
}

func (h *Handler) CreateProvince(ctx *fiber.Ctx) error {
	var payload ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	province, err := h.service.Create(payload)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("Province created", province))
}

func (h *Handler) UpdateProvince(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	province, err := h.service.Update(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Province updated", province))
}

func (h *Handler) DeleteProvince(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Province deleted", nil))
}
