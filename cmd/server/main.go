package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/redis"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/shutdown"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func main() {
	bootstrap.InitApp()

	app := fiber.New(fiber.Config{
		AppName:                 config.AppConfig.AppName,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
	})

	configureMiddleware(app)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient := redis.Connect(config.AppConfig)

	bootstrap.Initialize(
		app, db, redisClient.CacheClient, redisClient.LockClient,
	)

	app.Static("/files", "./storage/uploads", fiber.Static{
		Browse: false,
		MaxAge: 3600,
	})

	// Setup graceful shutdown
	shutdownHandler := shutdown.NewShutdownHandler(
		app,
		db,
		redisClient,
	).WithTimeout(20 * time.Second)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port)); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Not Found Resource Error Helper
	app.Use(helpers.NotFoundHelper)

	// Block main thread until shutdown signal
	shutdownHandler.ListenForShutdown()
}

func configureMiddleware(app *fiber.App) {
	// Initialize default config
	app.Use(logger.New())

	// Add Request ID middleware
	app.Use(requestid.New())

	app.Use(helpers.APILogger(helpers.GetAPILogger()))

	// Recover panic
	app.Use(helpers.RecoverWithLog())

	// Error helper
	app.Use(helpers.ErrorHelper)
}
