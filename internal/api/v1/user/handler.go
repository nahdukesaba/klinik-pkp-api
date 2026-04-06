package user

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUsersHandler(ctx *fiber.Ctx) error {
	pagination, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, err.Error(), err, true)
	}

	data, totalRecords, err := h.service.GetUsers(pagination)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", utils.NewPaginatedData(data, totalRecords, pagination), false)
}

func (h *Handler) GetUserByIdHandler(ctx *fiber.Ctx) error {
	var authID = ctx.Locals("id").(uuid.UUID)
	var userID uuid.UUID
	var err error

	if ctx.Params("id") == "me" {
		// SET USER ID TO AUTHENTICATED USER ID IF PARAMETER IS "me"
		userID = authID
	} else {
		// SET USER ID TO PARAMETER VALUE IF NOT "me"
		userID, err = uuid.Parse(ctx.Params("id"))

		if err != nil {
			return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid user ID", err, true)
		}

		// PREVENT ACCESS TO OTHER USERS' PROFILES EXCEPT FOR ADMINS
		if userID != authID && ctx.Locals("role").(string) != "admin" {
			return utils.JSONResponse(ctx, fiber.StatusForbidden, fiber.ErrForbidden.Message, fiber.ErrForbidden, true)
		}
	}

	data, err := h.service.GetUserById(userID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) PostUserHandler(ctx *fiber.Ctx) error {
	var payload AddUserPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddUser(&payload)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusCreated, "user created", fiber.Map{
		"id": data.ID,
	}, false)
}

// UPDATE PROFILE UPDATES CURRENT USER PROFILE
func (h *Handler) PutUserByIdHandler(ctx *fiber.Ctx) error {
	authID := ctx.Locals("id").(uuid.UUID)
	userID, err := uuid.Parse(ctx.Params("id"))

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid user ID", err, true)
	}

	if userID != authID {
		return utils.JSONResponse(ctx, fiber.StatusForbidden, fiber.ErrForbidden.Message, fiber.ErrForbidden, true)
	}

	var payload EditUserProfilePayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	if err := h.service.EditUserById(userID, &payload); err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	return utils.JSONResponse(ctx, fiber.StatusOK, "profile updated", nil, false)
}
