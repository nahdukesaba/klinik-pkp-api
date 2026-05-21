package main

import (
	"klinik-pkp-api/config"
	"klinik-pkp-api/internal/api/v1"
	"klinik-pkp-api/internal/middleware"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()

	// SET DATABASE CONNECTION
	db := config.ConnectDatabase(cfg.GetDSN())

	// INIT HTTP SERVER
	app := fiber.New(fiber.Config{
		AppName: "Klinik PKP Sumatera Utara API",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
		BodyLimit: 4 * 1024 * 1024, // 4MB max body size
	})

	// ENABLE CORS WITH CREDENTIALS FOR COOKIE-BASED REFRESH TOKEN FLOW
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORSAllowOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// /api/v1 ROUTES
	apiVersion1 := app.Group("/api/v1")
	v1.SetupRoutes(apiVersion1, db)

	// /api/v2 ROUTES (PLACEHOLDER FOR FUTURE VERSION)
	// apiVersion2 := app.Group("/api/v2")
	// v2.SetupRoutes(apiVersion2, db)

	middleware.StartViewerWorker(db)
	app.Use(middleware.ViewerMiddleware())

	log.Printf("Server starting on port %s (Environment: %s)", cfg.AppPort, cfg.AppEnv)

	if err := app.Listen(cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
