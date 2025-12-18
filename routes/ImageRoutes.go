package routes

import (
	"klinik-api/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupImageRoutes(router fiber.Router, ctrl *controller.ImageController) {
	images := router.Group("/images")

	// Public routes - baca image
	images.Get("/", ctrl.GetAllImages)
	images.Get("/:id", ctrl.GetImageByID)
	images.Get("/entity", ctrl.GetImagesByEntity)

	// Protected routes - upload dan hapus (butuh authentication)
	// Uncomment jika sudah ada middleware auth
	// images.Post("/upload", middleware.Auth(), ctrl.UploadImage)
	// images.Post("/upload-multiple", middleware.Auth(), ctrl.UploadMultipleImages)
	// images.Put("/:id", middleware.Auth(), ctrl.UpdateImageEntity)
	// images.Delete("/:id", middleware.Auth(), ctrl.DeleteImage)

	// Sementara tanpa auth (development only)
	images.Post("/upload", ctrl.UploadImage)
	images.Post("/upload-multiple", ctrl.UploadMultipleImages)
	images.Put("/:id", ctrl.UpdateImageEntity)
	images.Delete("/:id", ctrl.DeleteImage)
}
