package bank_desain

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	bankDesain := app.Group("/bank-desain")

	bankDesain.Get("/", handler.GetBankDesainHandler)
	bankDesain.Get("/:id", handler.GetBankDesainByIdHandler)
	bankDesain.Post("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostBankDesainHandler)
	bankDesain.Put("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutBankDesainByIdHandler)
	bankDesain.Delete("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteBankDesainByIdHandler)
}
