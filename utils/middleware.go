package utils

import (
	"errors"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AUTH MIDDLEWARE - ACCEPTS LEGACY KLINIK JWT OR SSO TOKEN.
func AuthMiddleware() fiber.Handler {
	tokenManager := NewTokenManager()
	ssoClient := NewSSOClient()

	return func(ctx *fiber.Ctx) error {
		authHeader := strings.TrimSpace(
			ctx.Get("Authorization"),
		)

		if authHeader == "" {
			return JSONResponse(
				ctx,
				fiber.StatusUnauthorized,
				fiber.ErrUnauthorized.Message,
				fiber.ErrUnauthorized,
				true,
			)
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {

			return JSONResponse(
				ctx,
				fiber.StatusUnauthorized,
				fiber.ErrUnauthorized.Message,
				fiber.ErrUnauthorized,
				true,
			)
		}

		accessToken := parts[1]

		// FIRST: preserve compatibility with existing Klinik JWT.
		claims, err := tokenManager.VerifyAccessToken(
			accessToken,
		)
		if err == nil {
			ctx.Locals("id", claims.UserID)
			ctx.Locals("email", claims.Email)
			ctx.Locals("role", claims.Role)

			return ctx.Next()
		}

		// SECOND: fallback to SSO.
		ssoContext, ssoErr := ssoClient.ResolveKlinikContext(
			ctx.UserContext(),
			accessToken,
		)
		if ssoErr != nil {
			if errors.Is(
				ssoErr,
				ErrSSOUnavailable,
			) {
				return JSONResponse(
					ctx,
					fiber.StatusServiceUnavailable,
					"SSO service unavailable",
					fiber.ErrServiceUnavailable,
					true,
				)
			}

			return JSONResponse(
				ctx,
				fiber.StatusUnauthorized,
				fiber.ErrUnauthorized.Message,
				fiber.ErrUnauthorized,
				true,
			)
		}

		ctx.Locals("id", ssoContext.ID)
		ctx.Locals("email", ssoContext.Email)
		ctx.Locals("role", ssoContext.Role)

		return ctx.Next()
	}
}

// ROLE MIDDLEWARE - RESTRICTS ACCESS TO SPECIFIC ROLES.
// MUST BE USED AFTER AUTH MIDDLEWARE.
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role, ok := ctx.Locals("role").(string)
		if !ok {
			return JSONResponse(
				ctx,
				fiber.StatusForbidden,
				fiber.ErrForbidden.Message,
				fiber.ErrForbidden,
				true,
			)
		}

		if slices.Contains(allowedRoles, role) {
			return ctx.Next()
		}

		return JSONResponse(
			ctx,
			fiber.StatusForbidden,
			fiber.ErrForbidden.Message,
			fiber.ErrForbidden,
			true,
		)
	}
}
