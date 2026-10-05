package sosialisasi

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/utils"
)

func SetupRoutes(app fiber.Router, handler *Handler) {
	sosialisasi := app.Group("/sosialisasi")

	sosialisasi.Get("/", handler.GetSosialisasiHandler)
	sosialisasi.Get("/:id", handler.GetSosialisasiByIdHandler)
	sosialisasi.Post("/", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PostSosialisasiHandler)
	sosialisasi.Put("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.PutSosialisasiByIdHandler)
	sosialisasi.Delete("/:id", utils.AuthMiddleware(), utils.RoleMiddleware("admin"), handler.DeleteSosialisasiByIdHandler)
}
