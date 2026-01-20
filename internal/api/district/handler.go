package district

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	response, err := h.service.GetDistricts()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", response, false)
}

func (h *Handler) GetDistrictByIDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	data, err := h.service.GetDistrictById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "district not found", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) PostDistrictHandler(ctx *fiber.Ctx) error {
	var payload PostDistrictPayload
	validator := utils.NewValidator()

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if err := validator.ValidateStruct(payload); err != nil {
		message := utils.FormatValidationError(err)

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddDistrict(payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "district created", data, false)
}

func (h *Handler) PutDistrictByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload PutDistrictPayload
	validator := utils.NewValidator()

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if err := validator.ValidateStruct(payload); err != nil {
		message := utils.FormatValidationError(err)

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditDistrictById(id, payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "district not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "district updated", nil, false)
}

func (h *Handler) DeleteDistrictByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteDistrictById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "district not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "district deleted", nil, false)
}
