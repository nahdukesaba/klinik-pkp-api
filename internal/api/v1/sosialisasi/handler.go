package sosialisasi

import (
	"encoding/json"
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

func (h *Handler) GetSosialisasiHandler(ctx *fiber.Ctx) error {
	var filter SosialisasiFilter

	// PARSE OPTIONAL QUERY PARAMETERS
	if val := ctx.Query("village_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid village_id", err, true)
		}

		filter.VillageID = &id
	}

	if val := ctx.Query("district_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid district_id", err, true)
		}

		filter.DistrictID = &id
	}

	if val := ctx.Query("region_id"); val != "" {
		id, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid region_id", err, true)
		}

		filter.RegionID = &id
	}

	if val := ctx.Query("location"); val != "" {
		filter.Location = &val
	}

	if val := ctx.Query("title"); val != "" {
		filter.Title = &val
	}

	data, err := h.service.GetSosialisasi(filter)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetSosialisasiV1Responses(data), false)
}

func (h *Handler) GetSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid sosialisasi id", err, true)
	}

	data, err := h.service.GetSosialisasiById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetSosialisasiV1Response(*data), false)
}

func (h *Handler) PostSosialisasiHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid multipart form", err, true)
	}

	var payload SosialisasiPayload

	// PARSE MULTIPART FORM
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// ASSIGN IMAGES FROM MULTIPART FORM
	payload.Images = form.File["images"]

	// PARSE COORDINATES FROM RAW STRING
	if err := json.Unmarshal([]byte(payload.CoordinateRaw), &payload.Coordinate); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid coordinates format", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddSosialisasi(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "sosialisasi created", AddSosialisasiV1Response{
		ID:        fmt.Sprintf("%d", data.ID),
		ImageURLs: data.ImageURLs,
	}, false)
}

func (h *Handler) PutSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid sosialisasi id", err, true)
	}

	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	var payload SosialisasiPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	payload.Images = form.File["images"]

	// PARSE COORDINATES FROM RAW STRING
	if err := json.Unmarshal([]byte(payload.CoordinateRaw), &payload.Coordinate); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid coordinates format", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditSosialisasiById(id, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "sosialisasi updated", nil, false)
}

func (h *Handler) DeleteSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid sosialisasi id", err, true)
	}

	if err := h.service.DeleteSosialisasiById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "sosialisasi deleted", nil, false)
}
