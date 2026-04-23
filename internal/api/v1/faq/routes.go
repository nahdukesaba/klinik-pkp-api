package faq

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app fiber.Router, handler *Handler) {
	app.Get("/faqs", handler.GetFAQsHandler)
	app.Get("/faqs/:id", handler.GetFAQByIdHandler)
	app.Post("/faqs", handler.PostFAQHandler)
	app.Put("/faqs/:id", handler.PutFAQHandler)
	app.Delete("/faqs/:id", handler.DeleteFAQHandler)
}
