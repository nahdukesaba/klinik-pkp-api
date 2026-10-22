package main

import (
	"context"
	"klinik-pkp-api/config"
	v1 "klinik-pkp-api/internal/api/v1"
	"klinik-pkp-api/utils"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// LOAD CONFIGURATION
	cfg := config.LoadConfig()
	if err := cfg.ValidateBrowserSSO(); err != nil {
		log.Fatalf("invalid browser SSO configuration: %v", err)
	}

	// FAIL FAST IF JWT CONFIGURATION IS UNSAFE
	if err := utils.ValidateJWTConfig(); err != nil {
		log.Fatalf("invalid JWT configuration: %v", err)
	}

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

			message := err.Error()

			if code >= fiber.StatusInternalServerError {
				log.Printf("internal server error: method=%s path=%s error=%v", ctx.Method(), ctx.Path(), err)

				message = fiber.ErrInternalServerError.Message
			}

			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   message,
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
	v1.SetupRoutes(apiVersion1, db, cfg)

	// /api/v2 ROUTES (PLACEHOLDER FOR FUTURE VERSION)
	// apiVersion2 := app.Group("/api/v2")
	// v2.SetupRoutes(apiVersion2, db)

	log.Printf("Server starting on port %s (Environment: %s)", cfg.AppPort, cfg.AppEnv)

	go func() {
		if err := app.Listen(cfg.AppPort); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}
