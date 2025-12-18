package main

import (
	"klinik-api/config"
	"klinik-api/controller"
	"klinik-api/middleware"
	"klinik-api/models"
	"klinik-api/routes"
	"klinik-api/seed"
	"klinik-api/service"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// customErrorHandler handles Fiber errors globally
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
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db := config.ConnectDatabase(cfg.GetDSN())
	log.Println("Database connected successfully")

	// Run migrations
	models.RunMigrations(db)

	// Run seeder if DB_SEED=true or "seed" argument provided
	if cfg.DBSeed || (len(os.Args) > 1 && os.Args[1] == "seed") {
		seed.Run(db, &cfg)
		if len(os.Args) > 1 && os.Args[1] == "seed" {
			log.Println("Seeding completed! Check your database.")
			return
		}
	}

	// Initialize services
	userService := service.NewUserService(db)
	adminService := service.NewAdminService(db, &cfg)
	provinceService := service.NewProvinceService(db)
	regencyService := service.NewRegencyService(db)
	categoryService := service.NewCategoryService(db)
	balaiService := service.NewBalaiService(db)
	rusunService := service.NewRusunService(db)
	imageService := service.NewImageService(db)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	adminController := controller.NewAdminController(adminService)
	provinceController := controller.NewProvinceController(provinceService)
	regencyController := controller.NewRegencyController(regencyService)
	categoryController := controller.NewCategoryController(categoryService)
	balaiController := controller.NewBalaiController(balaiService)
	rusunController := controller.NewRusunController(rusunService)
	imageController := controller.NewImageController(imageService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Klinik PKP Sumatera Utara API",
		ErrorHandler: customErrorHandler,
		BodyLimit:    4 * 1024 * 1024, // 4MB max body size
	})

	// Global middleware (security first)
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())
	app.Use(middleware.Logger())
	app.Use(middleware.RateLimiter())

	// API routes
	api := app.Group("/api/v1")

	// Public authentication routes (with stricter rate limiting)
	authGroup := api.Group("/auth")
	authGroup.Use(middleware.AuthRateLimiter())
	routes.SetupAuthRoutes(authGroup, userController, adminController)

	// Protected routes
	routes.SetupUserRoutes(api, userController)
	routes.SetupAdminRoutes(api, adminController)
	
	// Public resource routes
	routes.SetupProvinceRoutes(api, provinceController)
	routes.SetupRegencyRoutes(api, regencyController)
	routes.SetupCategoryRoutes(api, categoryController)
	routes.SetupBalaiRoutes(api, balaiController)
	routes.SetupRusunRoutes(api, rusunController)
	routes.SetupImageRoutes(api, imageController)

	// Serve static files untuk uploads
	app.Static("/uploads", "./uploads")

	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Klinik PKP Sumatera Utara API",
			"version": "3.0.0",
			"docs":    "/api/docs",
			"env":     cfg.AppEnv,
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"env":      cfg.AppEnv,
		})
	})

	// 404 handler (must be last)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
			"path":    c.Path(),
		})
	})

	// Start server
	port := ":3000"
	log.Printf("Server starting on port %s (Environment: %s)", port, cfg.AppEnv)
	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

