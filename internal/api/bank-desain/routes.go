package bank_desain

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/bank-desain", handler.GetBankDesainHandler)
	app.Get("/bank-desain/:id", handler.GetBankDesainByIdHandler)
	app.Post("/bank-desain", handler.PostBankDesainHandler)
	app.Put("/bank-desain/:id", handler.PutBankDesainByIdHandler)
	app.Delete("/bank-desain/:id", handler.DeleteBankDesainByIdHandler)
}