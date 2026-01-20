package rusun

import (
	"klinik-pkp-api/utils"
	"strconv"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
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
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	rusun, err := h.service.GetRusunById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", rusun, false)
}

func (h *Handler) PostRusunHandler(ctx *fiber.Ctx) error {
	var payload RusunPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	rusun, err := h.service.AddRusun(payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "rusun created", rusun, false)
}

func (h *Handler) PutRusunByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	var payload RusunPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	rusun, err := h.service.EditRusunById(id, payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "Rusun updated", rusun, false)
}

func (h *Handler) DeleteRusunByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	err = h.service.DeleteRusunById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "rusun deleted", nil, false)
}
