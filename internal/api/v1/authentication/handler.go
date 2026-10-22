package authentication

import (
	"errors"

	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Handler struct {
	service    *Service
	browserSSO *BrowserSSOService
}

// CONSTRUCTOR
func NewHandler(service *Service, browserSSO ...*BrowserSSOService) *Handler {
	handler := &Handler{service: service}
	if len(browserSSO) > 0 {
		handler.browserSSO = browserSSO[0]
	}
	return handler
}

type browserSSOExchangePayload struct {
	Code   string `json:"code"`
	State  string `json:"state"`
	Status string `json:"status"`
}

func (h *Handler) BrowserSSOStartHandler(ctx *fiber.Ctx) error {
	if h.browserSSO == nil {
		return browserSSOHandlerError(ctx, browserSSOUnavailable())
	}
	authorizeURL, err := h.browserSSO.Start(ctx.UserContext())
	if err != nil {
		return browserSSOHandlerError(ctx, err)
	}
	return ctx.Redirect(authorizeURL, fiber.StatusFound)
}

func (h *Handler) BrowserSSOExchangeHandler(ctx *fiber.Ctx) error {
	if h.browserSSO == nil {
		return browserSSOHandlerError(ctx, browserSSOUnavailable())
	}
	var payload browserSSOExchangePayload
	if err := ctx.BodyParser(&payload); err != nil {
		return browserSSOHandlerError(ctx, &BrowserSSOError{
			StatusCode: fiber.StatusBadRequest, Code: "INVALID_BODY", Message: "invalid request body",
		})
	}
	result, err := h.browserSSO.Complete(
		ctx.UserContext(), payload.State, payload.Code, payload.Status,
	)
	if err != nil {
		return browserSSOHandlerError(ctx, err)
	}
	if result == nil {
		return browserSSOHandlerError(ctx, browserSSOUnavailable())
	}
	if result.Session == nil {
		return utils.JSONResponse(ctx, fiber.StatusOK, "Klinik access is pending", fiber.Map{
			"status": result.Status,
		}, false)
	}
	return utils.JSONResponse(ctx, fiber.StatusOK, "browser SSO authentication completed", fiber.Map{
		"status":       result.Status,
		"accessToken":  result.Session.AccessToken,
		"refreshToken": result.Session.RefreshToken,
		"tokenType":    result.Session.TokenType,
		"expiresIn":    result.Session.ExpiresIn,
		"expiresAt":    result.Session.ExpiresAt,
		"user":         result.User,
	}, false)
}

func (h *Handler) BrowserSSOLogoutHandler(ctx *fiber.Ctx) error {
	if h.browserSSO == nil {
		return browserSSOHandlerError(ctx, browserSSOUnavailable())
	}
	logoutURL, err := h.browserSSO.LogoutURL()
	if err != nil {
		return browserSSOHandlerError(ctx, err)
	}
	return ctx.Redirect(logoutURL, fiber.StatusFound)
}

func browserSSOHandlerError(ctx *fiber.Ctx, err error) error {
	var ssoError *BrowserSSOError
	if !errors.As(err, &ssoError) {
		ssoError = browserSSOUnavailable()
	}
	return ctx.Status(ssoError.StatusCode).JSON(fiber.Map{
		"success": false,
		"error":   ssoError.Code,
		"message": ssoError.Message,
	})
}

func (h *Handler) PostAuthenticationHandler(ctx *fiber.Ctx) error {
	var payload AddAuthenticationPayload

	// VALIDATE CONTENT TYPE
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, fiber.ErrBadRequest.Message, err, true)
	}

	// VALIDATE PAYLOAD STRUCT
	if message, err := h.service.validator.ValidateStruct(payload); err != nil {
		return utils.JSONResponse(ctx, fiber.StatusBadRequest, message, err, true)
	}

	data, err := h.service.AddAuthentication(ctx.UserContext(), payload)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) ||
			errors.Is(err, utils.ErrKlinikAuthorizationDenied) {
			return utils.JSONResponse(ctx, fiber.StatusUnauthorized, "invalid credentials", fiber.ErrUnauthorized, true)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, fiber.ErrInternalServerError.Message, fiber.ErrInternalServerError, true)
	}

	h.service.tokenManager.SetRefreshTokenCookie(ctx, data["refresh_token"].(string))

	return utils.JSONResponse(ctx, fiber.StatusCreated, "authentication created", fiber.Map{
		"id":           data["id"],
		"access_token": data["access_token"],
	}, false)
}

// VALIDATES REFRESH TOKEN, GENERATES NEW ACCESS TOKEN
func (h *Handler) PutAuthenticationHandler(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	if refreshToken == "" {
		return utils.JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
	}

	data, err := h.service.UpdateAuthentication(ctx.UserContext(), refreshToken)

	if err != nil {
		return utils.JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, err, true)
	}

	newRefreshToken, ok := data["refresh_token"].(string)
	if !ok || newRefreshToken == "" {
		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, fiber.ErrInternalServerError.Message, fiber.ErrInternalServerError, true)
	}

	h.service.tokenManager.SetRefreshTokenCookie(ctx, newRefreshToken)

	return utils.JSONResponse(
		ctx,
		fiber.StatusOK,
		"access token refreshed",
		fiber.Map{
			"access_token": data["access_token"],
		},
		false,
	)
}

// VALIDATES AND REMOVES REFRESH TOKEN FROM DATABASE
func (h *Handler) DeleteAuthenticationHandler(ctx *fiber.Ctx) error {
	centralLogoutRequired := false
	if trusted, ok := utils.TrustedAuthContext(ctx); ok {
		centralLogoutRequired = trusted.LoginSource == utils.LoginSourceSSO
	}
	responseData := fiber.Map{"centralLogoutRequired": centralLogoutRequired}
	refreshToken := ctx.Cookies("refresh_token")

	if refreshToken == "" {
		h.service.tokenManager.ClearRefreshTokenCookie(ctx)
		return utils.JSONResponse(ctx, fiber.StatusOK, "session terminated", responseData, false)
	}

	if err := h.service.DeleteAuthentication(ctx.UserContext(), refreshToken); err != nil {
		if err == gorm.ErrRecordNotFound {
			h.service.tokenManager.ClearRefreshTokenCookie(ctx)
			return utils.JSONResponse(ctx, fiber.StatusOK, "session terminated", responseData, false)
		}

		return utils.JSONResponse(ctx, fiber.StatusInternalServerError, "", err, true)
	}

	h.service.tokenManager.ClearRefreshTokenCookie(ctx)

	return utils.JSONResponse(ctx, fiber.StatusOK, "session terminated", responseData, false)
}
