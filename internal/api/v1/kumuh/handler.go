package kumuh

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

func (h *Handler) GetKumuhHandler(ctx *fiber.Ctx) error {
	var filter KawasanKumuhFilter

	// PARSE OPTIONAL QUERY PARAMETERS
	if val := ctx.Query("region_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid region_id", err, true)
		}

		filter.RegionID = &id
	}

	if val := ctx.Query("district_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid district_id", err, true)
		}

		filter.DistrictID = &id
	}

	if val := ctx.Query("village_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid village_id", err, true)
		}

		filter.VillageID = &id
	}

	if val := ctx.Query("area_name"); val != "" {
		filter.AreaName = &val
	}

	if val := ctx.Query("year_inspected"); val != "" {
		year, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid year_inspected", err, true)
		}

		filter.YearInspected = &year
	}

	data, err := h.service.GetKumuh(filter)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetKumuhV1Responses(data), false)
}

func (h *Handler) GetKumuhByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid kawasan kumuh id", err, true)
	}

	data, err := h.service.GetKumuhById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetKumuhV1Response(*data), false)
}

func (h *Handler) PostKumuhHandler(ctx *fiber.Ctx) error {
	var payload KawasanKumuhPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadGateway.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddKumuh(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "kawasan kumuh created", AddKumuhV1Response{
		ID: fmt.Sprintf("%d", data.ID),
	}, false)
}

func (h *Handler) PutKumuhByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid kawasan kumuh id", err, true)
	}

	var payload KawasanKumuhPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditKumuhById(id, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "kawasan kumuh updated", nil, false)
}

func (h *Handler) DeleteKumuhByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid kawasan kumuh id", err, true)
	}

	if err := h.service.DeleteKumuhById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "kawasan kumuh deleted", nil, false)
}
