package v1

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/v1/authentication"
	"klinik-pkp-api/internal/api/v1/bank-desain"
	"klinik-pkp-api/internal/api/v1/bsps"
	"klinik-pkp-api/internal/api/v1/district"
	"klinik-pkp-api/internal/api/v1/kumuh"
	"klinik-pkp-api/internal/api/v1/province"
	"klinik-pkp-api/internal/api/v1/region"
	"klinik-pkp-api/internal/api/v1/rusun"
	"klinik-pkp-api/internal/api/v1/sosialisasi"
	"klinik-pkp-api/internal/api/v1/uploads"
	"klinik-pkp-api/internal/api/v1/user"
	"klinik-pkp-api/internal/api/v1/village"
)

// REGISTER ALL ROUTES FOR API V1
func SetupRoutes(api fiber.Router, db *gorm.DB) {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// MAIN ROUTES
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"docs":    "/api/docs",
			"env":     cfg.AppEnv,
			"message": "Welcome to Klinik PKP Sumatera II API",
			"version": "v1",
		})
	})
	api.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"env":      cfg.AppEnv,
		})
	})

	// PUBLIC ROUTES
	authentication.SetupRoutes(api, authentication.NewHandler(authentication.NewService(db)))
	user.SetupRoutes(api, user.NewHandler(user.NewService(db)))

	// DATA ROUTES
	province.SetupRoutes(api, province.NewHandler(province.NewService(db)))
	region.SetupRoutes(api, region.NewHandler(region.NewService(db)))
	district.SetupRoutes(api, district.NewHandler(district.NewService(db)))
	village.SetupRoutes(api, village.NewHandler(village.NewService(db)))
	bsps.SetupRoutes(api, bsps.NewHandler(bsps.NewService(db)))
	kumuh.SetupRoutes(api, kumuh.NewHandler(kumuh.NewService(db)))
	rusun.SetupRoutes(api, rusun.NewHandler(rusun.NewService(db, uploads.NewService("./storage"))))
	bank_desain.SetupRoutes(api, bank_desain.NewHandler(bank_desain.NewService(db, uploads.NewService("./storage"))))
	sosialisasi.SetupRoutes(api, sosialisasi.NewHandler(sosialisasi.NewService(db, uploads.NewService("./storage"))))
	uploads.SetupRoutes(api)

	// NOT FOUND ROUTE
	api.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   fiber.ErrNotFound.Message,
			"path":    ctx.Path(),
		})
	})
}
