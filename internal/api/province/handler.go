package province

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

func (h *Handler) GetProvincesHandler(ctx *fiber.Ctx) error {
	response, err := h.service.GetProvinces()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("success", response))
}

func (h *Handler) GetProvinceByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	response, err := h.service.GetProvinceById(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", response))
}

func (h *Handler) PostProvinceHandler(ctx *fiber.Ctx) error {
	var payload ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	response, err := h.service.AddProvince(payload)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}
	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("province created", response))
}

func (h *Handler) PutProvinceByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload ProvincePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	response, err := h.service.EditProvinceById(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Province updated", response))
}

func (h *Handler) DeleteProvinceByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteProvinceById(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("province deleted", nil))
}
