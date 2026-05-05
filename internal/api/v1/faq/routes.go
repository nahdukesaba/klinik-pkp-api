package faq

import (
	"klinik-pkp-api/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/faqs", handler.GetFAQsHandler)
	app.Get("/faqs/:id", handler.GetFAQByIdHandler)
	app.Post("/faqs", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostFAQHandler)
	app.Put("/faqs/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutFAQHandler)
	app.Delete("/faqs/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteFAQHandler)
}
