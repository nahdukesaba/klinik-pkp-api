package sosialisasi

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"klinik-pkp-api/utils"
	"strconv"
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
	data, err := h.service.GetSosialisasi()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) GetSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid sosialisasi id", err, true)
	}

	data, err := h.service.GetSosialisasiById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "sosialisasi not found", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) PostSosialisasiHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid multipart form", err, true)
	}

	var payload SosialisasiPayload
	validator := utils.NewValidator()

	// PARSE MULTIPART FORM
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid payload", err, true)
	}

	// ASSIGN IMAGES FROM MULTIPART FORM
	payload.Images = form.File["images"]

	// PARSE COORDINATES FROM RAW STRING
	if err := json.Unmarshal([]byte(payload.CoordinatesRaw), &payload.Coordinates); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid coordinates format", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddSosialisasi(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "sosialisasi created", fiber.Map{
		"id": fmt.Sprintf("%d", data.ID),
		"image_urls": data.ImageURLs,
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
	validator := utils.NewValidator()

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	payload.Images = form.File["images"]

	// PARSE COORDINATES FROM RAW STRING
	if err := json.Unmarshal([]byte(payload.CoordinatesRaw), &payload.Coordinates); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid coordinates format", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditSosialisasiById(id, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "sosialisasi not found", err, true)
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
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "sosialisasi not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "sosialisasi deleted", nil, false)
}
