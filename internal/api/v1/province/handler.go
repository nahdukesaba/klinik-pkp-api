package province

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

func (h *Handler) GetProvincesHandler(ctx *fiber.Ctx) error {
	data, err := h.service.GetProvinces()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetProvinceV1Responses(data), false)
}

func (h *Handler) GetProvinceByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid province id", err, true)
	}

	data, err := h.service.GetProvinceById(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", ToGetProvinceV1Response(*data), false)
}

func (h *Handler) PostProvinceHandler(ctx *fiber.Ctx) error {
	var payload ProvincePayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid content-type", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddProvince(payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "province created", AddProvinceV1Response{
		ID: fmt.Sprintf("%d", data.ID),
	}, false)
}

func (h *Handler) PutProvinceByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid province id", err, true)
	}

	var payload ProvincePayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid content-type", err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditProvinceById(id, payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "province updated", nil, false)
}

func (h *Handler) DeleteProvinceByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid province id", err, true)
	}

	if err := h.service.DeleteProvinceById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "province deleted", nil, false)
}
