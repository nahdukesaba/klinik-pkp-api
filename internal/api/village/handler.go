package village

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

func (h *Handler) GetVillagesHandler(ctx *fiber.Ctx) error {
	villages, err := h.service.GetVillages()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.ErrBadRequest))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", villages))
}

func (h *Handler) GetVillageByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	village, err := h.service.GetVillageById(id)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Success", village))
}

func (h *Handler) PostVillageHandler(ctx *fiber.Ctx) error {
	var payload VillagePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	village, err := h.service.AddVillage(payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("Village created", village))
}

func (h *Handler) PutVillageByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var payload VillagePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	village, err := h.service.EditVillageById(id, payload)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Village updated", village))
}

func (h *Handler) DeleteVillageByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.DeleteVillageById(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse("Village deleted", nil))
}
