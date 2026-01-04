package main

import (
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/province"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/rusun"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/internal/api/bsps"
	"klinik-pkp-api/internal/api/user"
	"log"
	"github.com/gofiber/fiber/v2"
)

// HANDLES FIBER ERRORS GLOBALLY
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	// INIT SERVICE FOR EACH MODULE
	userService := user.NewService(db)
	rusunService := rusun.NewService(db)
	provinceService := province.NewService(db)
	regionService := region.NewService(db)
	districtService := district.NewService(db)
	villageService := village.NewService(db)
	bspsService := bsps.NewService(db)

	// INIT HANDLER FOR EACH SERVICE
	userHandler := user.NewHandler(userService)
	rusunHandler := rusun.NewHandler(rusunService)
	provinceHandler := province.NewHandler(provinceService)
	regionHandler := region.NewHandler(regionService)
	districtHandler := district.NewHandler(districtService)
	villageHandler := village.NewHandler(villageService)
	bspsHandler := bsps.NewHandler(bspsService)

	// INIT HTTP SERVER
	app := fiber.New(fiber.Config{
		AppName:      "Klinik PKP Sumatera Utara API",
		ErrorHandler: customErrorHandler,
		BodyLimit:    4 * 1024 * 1024, // 4MB max body size
	})

	// SERVE STATIC FILES FOR UPLOADS
	app.Static("/uploads", "./uploads")

	// GROUP ROUTES WITH /api/v1 PREFIX
	api := app.Group("/api/v1")

	// ALL ROUTES
	user.SetupRoutes(api, userHandler)
	rusun.SetupRoutes(api, rusunHandler)
	province.SetupRoutes(api, provinceHandler)
	region.SetupRoutes(api, regionHandler)
	district.SetupRoutes(api, districtHandler)
	village.SetupRoutes(api, villageHandler)
	bsps.SetupRoutes(api, bspsHandler)

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
