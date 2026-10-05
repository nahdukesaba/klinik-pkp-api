package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// API RESPONSE STRUCTURE
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SUCCESS RESPONSE
func SuccessResponse(message string, data any) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ERROR RESPONSE
func ErrorResponse(message string, err error) Response {
	errorMessage := err.Error()

	if message != "" {
		errorMessage = message
	}

	return Response{
		Success: false,
		Error:   errorMessage,
	}
}

// GENERIC JSON RESPONSE
func JSONResponse(ctx *fiber.Ctx, statusCode int, message string, data any, isError bool) error {
	if isError {
		err, ok := data.(error)

		if !ok {
			err = fiber.NewError(statusCode, message)
		}

		// NEVER EXPOSE INTERNAL SERVER ERRORS TO CLIENTS
		if statusCode >= fiber.StatusInternalServerError {
			log.Printf(
				"internal server error: method=%s path=%s error=%v",
				ctx.Method(),
				ctx.Path(),
				err,
			)

			return ctx.Status(statusCode).JSON(
				ErrorResponse(
					fiber.ErrInternalServerError.Message,
					fiber.ErrInternalServerError,
				),
			)
		}

		return ctx.Status(statusCode).JSON(ErrorResponse(message, err),)
	}

	return ctx.Status(statusCode).JSON(SuccessResponse(message, data),)
}
