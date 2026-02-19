package rusun

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

func (h *Handler) GetRusunHandler(ctx *fiber.Ctx) error {
	rusun, err := h.service.GetRusun()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", rusun, false)
}

func (h *Handler) GetRusunByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid rusun id", err, true)
	}

	rusun, err := h.service.GetRusunById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "rusun not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", rusun, false)
}

func (h *Handler) PostRusunHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid multipart form", err, true)
	}

	var payload RusunPayload
	validator := utils.NewValidator()

	// PARSE MULTIPART FORM
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid payload", err, true)
	}

	// ASSIGN IMAGES FROM MULTIPART FORM
	payload.Images = form.File["images"]

	// PARSE COORDINATES FROM RAW STRING
	if err := json.Unmarshal([]byte(payload.CoordinateRaw), &payload.Coordinate); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid coordinate format", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddRusun(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "rusun created", fiber.Map{
		"id":         fmt.Sprintf("%d", data.ID),
		"image_urls": data.ImageURLs,
	}, false)
}

func (h *Handler) PutRusunByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid rusun id", err, true)
	}

	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid payload", err, true)
	}

	var payload RusunPayload
	validator := utils.NewValidator()

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid payload", err, true)
	}

	payload.Images = form.File["images"]

	// PARSE COORDINATES FROM RAW STRING
	if err := json.Unmarshal([]byte(payload.CoordinateRaw), &payload.Coordinate); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid coordinates format", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditRusunById(id, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "rusun not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "rusun updated", nil, false)
}

func (h *Handler) DeleteRusunByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid rusun id", err, true)
	}

	if err := h.service.DeleteRusunById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "rusun not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "rusun deleted", nil, false)
}
