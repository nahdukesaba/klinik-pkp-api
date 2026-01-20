package sosialisasi

import (
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
	sosialisasi, err := h.service.GetSosialisasi()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", sosialisasi, false)
}

func (h *Handler) GetSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))

	sosialisasi, err := h.service.GetSosialisasiById(id)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", sosialisasi, false)
}

func (h *Handler) AddSosialisasiHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	payload := &SosialisasiPayload{
		VillageID:        ctx.FormValue("village_id"),
		DistrictID:       ctx.FormValue("district_id"),
		RegionID:         ctx.FormValue("region_id"),
		Title:            ctx.FormValue("title"),
		Location:         ctx.FormValue("location"),
		Description:      ctx.FormValue("description"),
		ScheduledAtStart: ctx.FormValue("scheduled_at_start"),
		ScheduledAtEnd:   ctx.FormValue("scheduled_at_end"),
		Images:           form.File["images"],
	}

	sosialisasi, err := h.service.AddSosialisasi(payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "sosialisasi created", sosialisasi, false)
}

func (h *Handler) PutSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	// PARSE MULTIPART FORM
	form, err := ctx.MultipartForm()

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	payload := &SosialisasiPayload{
		VillageID:        ctx.FormValue("village_id"),
		DistrictID:       ctx.FormValue("district_id"),
		RegionID:         ctx.FormValue("region_id"),
		Title:            ctx.FormValue("title"),
		Location:         ctx.FormValue("location"),
		Description:      ctx.FormValue("description"),
		ScheduledAtStart: ctx.FormValue("scheduled_at_start"),
		ScheduledAtEnd:   ctx.FormValue("scheduled_at_end"),
		Images:           form.File["images"],
	}

	sosialisasi, err := h.service.EditSosialisasiById(id, payload)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "sosialisasi updated", sosialisasi, false)
}

func (h *Handler) DeleteSosialisasiByIdHandler(ctx *fiber.Ctx) error {
	// PARSE ID TO UINT64
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	if err := h.service.DeleteSosialisasiById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, "", err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "sosialisasi deleted", nil, false)
}
