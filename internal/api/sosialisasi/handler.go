package sosialisasi

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

func (h *Handler) GetSosialisasiHandler(ctx *fiber.Ctx) error {
	sosialisasi, err := h.service.GetSosialisasi()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("success", sosialisasi))
}