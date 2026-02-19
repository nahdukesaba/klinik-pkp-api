package main

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/bank-desain"
	"klinik-pkp-api/internal/api/bsps"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/kumuh"
	"klinik-pkp-api/internal/api/province"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/rusun"
	"klinik-pkp-api/internal/api/sosialisasi"
	"klinik-pkp-api/internal/api/uploads"
	"klinik-pkp-api/internal/api/user"
	"klinik-pkp-api/internal/api/village"
	"log"
)

// HANDLES FIBER ERRORS GLOBALLY
func customErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return ctx.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	// INIT HTTP SERVER
	app := fiber.New(fiber.Config{
		AppName:      "Klinik PKP Sumatera Utara API",
		ErrorHandler: customErrorHandler,
		BodyLimit:    4 * 1024 * 1024, // 4MB max body size
	})

	// GROUP ROUTES WITH /api/v1 PREFIX
	api := app.Group("/api/v1")

	// ALL ROUTES
	user.SetupRoutes(api, user.NewHandler(user.NewService(db)))
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

	// MAIN ROUTE
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"docs":    "/api/docs",
			"env":     cfg.AppEnv,
			"message": "Welcome to Klinik PKP Sumatera II API",
			"version": "1.0.0",
		})
	})

	// HEALTH CHECK ROUTE
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"env":      cfg.AppEnv,
		})
	})

	// NOT FOUND ROUTE
	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
			"path":    c.Path(),
		})
	})

	log.Printf("Server starting on port %s (Environment: %s)", cfg.AppPort, cfg.AppEnv)

	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
