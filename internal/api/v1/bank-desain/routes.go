package bank_desain

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	bankDesain := app.Group("/bank-desain")

	bankDesain.Get("/", handler.GetBankDesainHandler)
	bankDesain.Get("/:id", handler.GetBankDesainByIdHandler)
	bankDesain.Post("/", handler.PostBankDesainHandler)
	bankDesain.Put("/:id", handler.PutBankDesainByIdHandler)
	bankDesain.Delete("/:id", handler.DeleteBankDesainByIdHandler)
	bankDesain.Post("/:id/download", handler.DownloadBankDesainHandler)
}
