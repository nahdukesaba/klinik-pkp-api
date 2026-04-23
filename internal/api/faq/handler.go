package faq

import (
	"strconv"

	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetFAQsHandler(ctx *fiber.Ctx) error {
	// DEFAULTS
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	search := ctx.Query("search", "")

	data, total, err := h.service.GetFAQs(page, limit, search)

	if err != nil {
		return utils.JSONResponse(ctx, 500, "", err, true)
	}

	return utils.JSONResponse(ctx, 200, "success", fiber.Map{
		"data":       data,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"total_page": (int(total) + limit - 1) / limit,
	}, false)
}

func (h *Handler) GetFAQByIdHandler(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)

	data, err := h.service.GetFAQById(id)

	if err != nil {
		return utils.JSONResponse(ctx, 404, "", err, true)
	}

	return utils.JSONResponse(ctx, 200, "success", data, false)
}

func (h *Handler) PostFAQHandler(ctx *fiber.Ctx) error {
	payload := new(FAQPayload)

	if err := ctx.BodyParser(payload); err != nil {
		return utils.JSONResponse(ctx, 400, "invalid request", err, true)
	}

	data, err := h.service.CreateFAQ(payload)

	if err != nil {
		return utils.JSONResponse(ctx, 500, "", err, true)
	}

	return utils.JSONResponse(ctx, 201, "created", data, false)
}

func (h *Handler) PutFAQHandler(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)

	payload := new(FAQPayload)

	if err := ctx.BodyParser(payload); err != nil {
		return utils.JSONResponse(ctx, 400, "invalid request", err, true)
	}

	if err := h.service.UpdateFAQ(id, payload); err != nil {
		return utils.JSONResponse(ctx, 500, "", err, true)
	}

	return utils.JSONResponse(ctx, 200, "updated", nil, false)
}

func (h *Handler) DeleteFAQHandler(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err := h.service.DeleteFAQ(id); err != nil {
		return utils.JSONResponse(ctx, 500, "", err, true)
	}

	return utils.JSONResponse(ctx, 200, "deleted", nil, false)
}
