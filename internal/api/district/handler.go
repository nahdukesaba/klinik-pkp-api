package district

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

func (h *Handler) GetAllDistricts(ctx *fiber.Ctx) error {
	districts, err := h.service.GetAll()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to fetch districts"))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("success", districts))
}

func (h *Handler) GetDistrictByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	district, err := h.service.GetByID(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("Success", district))
}

func (h *Handler) CreateDistrict(ctx *fiber.Ctx) error {
	var payload DistrictPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	district, err := h.service.Create(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessMessageResponse("District created", district))
}

func (h *Handler) UpdateDistrict(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload DistrictPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	district, err := h.service.Update(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("District updated", district))
}

func (h *Handler) DeleteDistrict(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessMessageResponse("District deleted", nil))
}
