package middleware

import (
	"klinik-api/config"
	"klinik-api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token with enhanced security
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Response{
				Success: false,
				Error:   "Missing authorization header",
			})
		}

		// Extract token (format: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Response{
				Success: false,
				Error:   "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token length (prevent excessively long tokens)
		if len(token) > 1000 {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Response{
				Success: false,
				Error:   "Invalid token",
			})
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.Response{
				Success: false,
				Error:   "Invalid or expired token",
			})
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// RequireAdmin middleware ensures user has admin role
func RequireAdmin() fiber.Handler {
	cfg := config.LoadConfig()
	
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(utils.Response{
				Success: false,
				Error:   "Access denied: No role found",
			})
		}

		// Check if role is any admin role
		adminRoles := []string{
			cfg.RoleSuperAdmin,
			cfg.RoleAdminEselon1,
			cfg.RoleAdminBalai,
		}

		isAdmin := false
		for _, adminRole := range adminRoles {
			if role == adminRole {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(utils.Response{
				Success: false,
				Error:   "Access denied: Admin privileges required",
			})
		}

		return c.Next()
	}
}

// RequireSuperAdmin middleware ensures user has super admin role
func RequireSuperAdmin() fiber.Handler {
	cfg := config.LoadConfig()
	
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(utils.Response{
				Success: false,
				Error:   "Access denied: No role found",
			})
		}

		if role != cfg.RoleSuperAdmin {
			return c.Status(fiber.StatusForbidden).JSON(utils.Response{
				Success: false,
				Error:   "Access denied: Super admin privileges required",
			})
		}

		return c.Next()
	}
}

// RequireRole middleware ensures user has specific role
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(utils.Response{
				Success: false,
				Error:   "Access denied: No role found",
			})
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(utils.Response{
			Success: false,
			Error:   "Access denied: Insufficient privileges",
		})
	}
}

// RequireAuth combines authentication check
func RequireAuth() fiber.Handler {
	return AuthMiddleware()
}
