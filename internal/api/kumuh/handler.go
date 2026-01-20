package kumuh

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
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
	data, err := h.service.GetKumuh()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}
