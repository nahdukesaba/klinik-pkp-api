package bank_desain

import (
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

func (h *Handler) GetBankDesainHandler(ctx *fiber.Ctx) error {
	data, err := h.service.GetBankDesain()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) GetBankDesainByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid bank desain id", err, true)
	}

	data, err := h.service.GetBankDesainById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "bank desain not found", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) PostBankDesainHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid multipart form", err, true)
	}

	var payload BankDesainPayload
	validator := utils.NewValidator()

	// PARSE REQUEST BODY
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid payload", err, true)
	}

	// ASSIGN IMAGES AND FILES FROM MULTIPART FORM
	payload.Images = form.File["images"]
	payload.Files = form.File["files"]

	// VALIDATE PAYLOAD STRUCT
	if message, err := validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddBankDesain(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "bank desain created", fiber.Map{
		"id":         fmt.Sprintf("%d", data.ID),
		"image_urls": data.ImageURLs,
		"file_urls":  data.FileURLs,
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
	validator := utils.NewValidator()

	// PARSE REQUEST BODY
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid payload", err, true)
	}

	payload.Images = form.File["images"]
	payload.Files = form.File["files"]

	// VALIDATE PAYLOAD STRUCT
	if message, err := validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditBankDesainById(id, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "bank desain not found", err, true)
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
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "bank desain not found", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "bank desain deleted", nil, false)
}
