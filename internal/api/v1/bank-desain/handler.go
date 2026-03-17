package bank_desain

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

func (h *Handler) GetBankDesainHandler(ctx *fiber.Ctx) error {
	var filter BankDesainFilter

	// PARSE OPTIONAL QUERY PARAMETERS
	if val := ctx.Query("type"); val != "" {
		filter.Type = &val
	}

	if val := ctx.Query("bedroom_count"); val != "" {
		count, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid bedroom_count", err, true)
		}

		filter.BedroomCount = &count
	}

	if val := ctx.Query("bathroom_count"); val != "" {
		count, err := strconv.ParseUint(val, 10, 64)

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid bathroom_count", err, true)
		}

		filter.BathroomCount = &count
	}

	if val := ctx.Query("has_garage"); val != "" {
		hasGarage := val == "true" || val == "1"
		filter.HasGarage = &hasGarage
	}

	if val := ctx.Query("name"); val != "" {
		filter.Name = &val
	}

	data, err := h.service.GetBankDesain(filter)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetBankDesainV1Responses(data), false)
}

func (h *Handler) GetBankDesainByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid bank desain id", err, true)
	}

	data, err := h.service.GetBankDesainById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetBankDesainV1Response(*data), false)
}

func (h *Handler) PostBankDesainHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid multipart form", err, true)
	}

	var payload BankDesainPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// ASSIGN IMAGES AND FILES FROM MULTIPART FORM
	payload.Images = form.File["images"]
	payload.Files = form.File["files"]

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddBankDesain(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "bank desain created", AddBankDesainV1Response{
		ID:        fmt.Sprintf("%d", data.ID),
		ImageURLs: data.ImageURLs,
		FileURLs:  data.FileURLs,
	}, false)
}

func (h *Handler) PutBankDesainByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid bank desain id", err, true)
	}

	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	var payload BankDesainPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	payload.Images = form.File["images"]
	payload.Files = form.File["files"]

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditBankDesainById(id, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "bank desain updated", nil, false)
}

func (h *Handler) DeleteBankDesainByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid bank desain id", err, true)
	}

	if err := h.service.DeleteBankDesainById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "bank desain deleted", nil, false)
}
