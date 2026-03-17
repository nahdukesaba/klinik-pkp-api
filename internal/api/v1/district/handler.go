package district

import (
	"fmt"
	"klinik-pkp-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	var filter DistrictFilter

	// PARSE OPTIONAL QUERY PARAMETERS
	if val := ctx.Query("region_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid region_id", err, true)
		}

		filter.RegionID = &id
	}

	data, err := h.service.GetDistricts(filter)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetDistrictV1Responses(data), false)
}

func (h *Handler) GetDistrictByIDHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid district id", err, true)
	}

	data, err := h.service.GetDistrictById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetDistrictV1Response(*data), false)
}

func (h *Handler) PostDistrictHandler(ctx *fiber.Ctx) error {
	var payload DistrictPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddDistrict(payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "district created", AddDistrictV1Response{
		ID: fmt.Sprintf("%d", data.ID),
	}, false)
}

func (h *Handler) PutDistrictByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid district id", err, true)
	}

	var payload DistrictPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditDistrictById(id, payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "district updated", nil, false)
}

func (h *Handler) DeleteDistrictByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid village id", err, true)
	}

	if err := h.service.DeleteDistrictById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "district deleted", nil, false)
}
