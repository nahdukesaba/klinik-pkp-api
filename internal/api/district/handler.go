package district

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

func (h *Handler) GetDistrictsHandler(ctx *fiber.Ctx) error {
	districts, err := h.service.GetDistricts()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("success", districts))
}

func (h *Handler) GetDistrictByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	district, err := h.service.GetDistrictById(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", district))
}

func (h *Handler) PostDistrictHandler(ctx *fiber.Ctx) error {
	var payload DistrictPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	district, err := h.service.AddDistrict(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("District created", district))
}

func (h *Handler) PutDistrictByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload DistrictPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	district, err := h.service.EditDistrictById(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("District updated", district))
}

func (h *Handler) DeleteDistrictByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteDistrictById(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("District deleted", nil))
}
