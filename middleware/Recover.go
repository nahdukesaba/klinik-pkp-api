package middleware

import (
	"klinik-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Recover returns panic recovery middleware
func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			c.Status(fiber.StatusInternalServerError).JSON(utils.Response{
				Success: false,
				Error:   "Internal server error",
			})
		},
	})
}
