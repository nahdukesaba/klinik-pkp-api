package main

import (
	"klinik-pkp-api/config"
	v1 "klinik-pkp-api/internal/api/v1"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// ENABLE CORS WITH CREDENTIALS FOR COOKIE-BASED REFRESH TOKEN FLOW
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSAllowOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// GROUP ROUTES WITH /api/v1 PREFIX
	api := app.Group("/api/v1")
	v1.SetupRoutes(api, db)

	// MAIN ROUTE
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"docs":    "/api/docs",
			"env":     cfg.AppEnv,
			"message": "Welcome to Klinik PKP Sumatera II API",
			"version": "1.0.0",
		})
	})

	// HEALTH CHECK ROUTE
	api.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"env":      cfg.AppEnv,
		})
	})

	// NOT FOUND ROUTE
	api.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   fiber.ErrNotFound.Message,
			"path":    ctx.Path(),
		})
	})

	log.Printf("Server starting on port %s (Environment: %s)", cfg.AppPort, cfg.AppEnv)

	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
