package user

import (
	"errors"
	"strings"

	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) GetUsersHandler(ctx *fiber.Ctx) error {
	pagination, err := utils.ParsePaginationQuery(ctx)
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, err.Error(), err, true)
	}
	query := DirectoryQuery{Pagination: pagination, Keyword: ctx.Query("keyword"), Role: MembershipRole(ctx.Query("role")), Status: MembershipStatus(ctx.Query("status")), SortBy: ctx.Query("sort_by"), SortOrder: ctx.Query("sort_order")}
	data, total, err := h.service.List(ctx.UserContext(), query)
	if err != nil {
		return handleDirectoryError(ctx, err)
	}
	return utils.JSONResponse(ctx, fiber.StatusOK, "success", utils.NewPaginatedData(data, total, pagination), false)
}

func (h *Handler) GetCurrentUserHandler(ctx *fiber.Ctx) error {
	trusted, ok := utils.TrustedAuthContext(ctx)
	if !ok {
		return utils.JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
	}
	data, err := h.service.GetCurrentUser(ctx.UserContext(), trusted.AppUserID)
	if err != nil {
		if errors.Is(err, ErrCanonicalIdentityUnavailable) {
			return utils.JSONResponse(ctx, fiber.StatusNotFound, fiber.ErrNotFound.Message, err, true)
		}
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}
	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) GetUserByIdHandler(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid app user ID", err, true)
	}
	data, err := h.service.Get(ctx.UserContext(), id)
	if err != nil {
		return handleDirectoryError(ctx, err)
	}
	return utils.JSONResponse(ctx, fiber.StatusOK, "success", data, false)
}

func (h *Handler) PostUserHandler(ctx *fiber.Ctx) error {
	var payload DirectoryCreate
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}
	data, err := h.service.Create(ctx.UserContext(), strings.TrimSpace(ctx.Get(fiber.HeaderAuthorization)), payload)
	if err != nil {
		return handleDirectoryError(ctx, err)
	}
	return utils.JSONResponse(ctx, fiber.StatusCreated, "Klinik access provisioned", data, false)
}

func (h *Handler) PutUserByIdHandler(ctx *fiber.Ctx) error {
	trusted, ok := utils.TrustedAuthContext(ctx)
	if !ok {
		return utils.JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
	}
	target, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid app user ID", err, true)
	}
	var payload DirectoryUpdate
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}
	data, err := h.service.Update(ctx.UserContext(), trusted.AppUserID, target, payload)
	if err != nil {
		return handleDirectoryError(ctx, err)
	}
	return utils.JSONResponse(ctx, fiber.StatusOK, "Klinik user updated", data, false)
}

func (h *Handler) DeleteUserByIdHandler(ctx *fiber.Ctx) error {
	trusted, ok := utils.TrustedAuthContext(ctx)
	if !ok {
		return utils.JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
	}
	target, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid app user ID", err, true)
	}
	if err := h.service.Disable(ctx.UserContext(), trusted.AppUserID, target); err != nil {
		return handleDirectoryError(ctx, err)
	}
	return utils.JSONResponse(ctx, fiber.StatusOK, "Klinik access disabled", nil, false)
}

func handleDirectoryError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrDirectoryUserNotFound):
		return utils.JSONResponse(ctx, fiber.StatusNotFound, "canonical Klinik user not found", err, true)
	case errors.Is(err, ErrInvalidDirectoryMutation):
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, "invalid canonical Klinik user data", err, true)
	case errors.Is(err, ErrSelfAccessMutation), errors.Is(err, ErrLastAdminMutation):
		return utils.JSONResponse(ctx, fiber.StatusConflict, err.Error(), err, true)
	case errors.Is(err, ErrProvisioningUnavailable):
		return utils.JSONResponse(ctx, fiber.StatusServiceUnavailable, "Central SSO provisioning unavailable", err, true)
	default:
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}
}
