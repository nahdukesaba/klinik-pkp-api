package main

import (
	"github.com/gofiber/fiber/v2"
	"klinik-pkp-api/config"
	"klinik-pkp-api/controller"
	"klinik-pkp-api/models"
	"klinik-pkp-api/routes"
	"klinik-pkp-api/seeders"
	"klinik-pkp-api/service"
	"log"
	"os"
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

	// RUN MIGRATION IF "migrate" ARGUMENT PROVIDED
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if len(os.Args) > 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
			models.DisplayAvailableMigrations()
		} else {
			// PASS ALL MODELS AFTER "migrate" COMMAND
			models.Run(db, os.Args[2:]...)
		}

		// EXIT AFTER MIGRATION
		return
	}

	// RUN SEEDER IF "seed" ARGUMENT PROVIDED
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if len(os.Args) > 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
			seeders.DisplayAvailableSeeders()
		} else {
			// PASS ALL SEEDS AFTER "seed" COMMAND
			seeders.Run(db, &cfg, os.Args[2:]...)
		}

		// EXIT AFTER SEEDING
		return
	}

	// INIT SERVICE FOR EACH MODELS
	userService := service.NewUserService(db)
	rusunService := service.NewRusunService(db)
	provinceService := service.NewProvinceService(db)
	regionService := service.NewRegionService(db)
	districtService := service.NewDistrictService(db)
	villageService := service.NewVillageService(db)

	// INIT CONTROLLER FOR EACH SERVICE
	userController := controller.NewUserController(userService)
	rusunController := controller.NewRusunController(rusunService)
	provinceController := controller.NewProvinceController(provinceService)
	regionController := controller.NewRegionController(regionService)
	districtController := controller.NewDistrictController(districtService)
	villageController := controller.NewVillageController(villageService)

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
	routes.SetupUserRoutes(api, userController)
	routes.SetupRusunRoutes(api, rusunController)
	routes.SetupProvinceRoutes(api, provinceController)
	routes.SetupRegionRoutes(api, regionController)
	routes.SetupDistrictRoutes(api, districtController)
	routes.SetupVillageRoutes(api, villageController)

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
