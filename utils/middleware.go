package utils

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AUTH MIDDLEWARE - VALIDATES ACCESS TOKEN ON PROTECTED ROUTES
func AuthMiddleware() fiber.Handler {
	tokenManager := NewTokenManager()

	return func(ctx *fiber.Ctx) error {
		// GET AUTHORIZATION HEADER
		authHeader := ctx.Get("Authorization")

		if authHeader == "" {
			return JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
		}

		// EXTRACT TOKEN FROM "Bearer <token>" FORMAT
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
		}

		// VERIFY ACCESS TOKEN SIGNATURE AND EXPIRATION
		claims, err := tokenManager.VerifyAccessToken(parts[1])
		if err != nil {
			return JSONResponse(ctx, fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message, fiber.ErrUnauthorized, true)
		}

		// STORE USER INFO IN CONTEXT LOCALS FOR NEXT HANDLERS
		ctx.Locals("id", claims.UserID)
		ctx.Locals("email", claims.Email)
		ctx.Locals("role", claims.Role)

		return ctx.Next()
	}
}

// ROLE MIDDLEWARE - RESTRICTS ACCESS TO SPECIFIC ROLES
// MUST BE USED AFTER AUTH MIDDLEWARE
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role, ok := ctx.Locals("role").(string)
		if !ok {
			return JSONResponse(ctx, fiber.StatusForbidden, fiber.ErrForbidden.Message, fiber.ErrForbidden, true)
		}

		// CHECK IF USER ROLE IS IN ALLOWED ROLES
		if slices.Contains(allowedRoles, role) {
			return ctx.Next()
		}

		return JSONResponse(ctx, fiber.StatusForbidden, fiber.ErrForbidden.Message, fiber.ErrForbidden, true)
	}
}
