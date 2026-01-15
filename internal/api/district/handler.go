package district

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

func (h *Handler) GetDistrictsHandler(ctx *fiber.Ctx) error {
	response, err := h.service.GetDistricts()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("success", response))
}

func (h *Handler) GetDistrictByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	response, err := h.service.GetDistrictById(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("success", response))
}

func (h *Handler) PostDistrictHandler(ctx *fiber.Ctx) error {
	var payload DistrictPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	response, err := h.service.AddDistrict(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("district created", response))
}

func (h *Handler) PutDistrictByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload DistrictPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	response, err := h.service.EditDistrictById(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("district updated", response))
}

func (h *Handler) DeleteDistrictByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteDistrictById(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("district deleted", nil))
}
