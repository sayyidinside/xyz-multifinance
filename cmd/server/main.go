package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func main() {
	bootstrap.InitApp()

	app := fiber.New(fiber.Config{
		AppName:                 config.AppConfig.AppName,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
	})

	// Initialize default config
	app.Use(logger.New())

	// Add Request ID middleware
	app.Use(requestid.New())

	app.Use(helpers.APILogger(helpers.GetAPILogger()))

	// Recover panic
	app.Use(helpers.RecoverWithLog())

	app.Use(helpers.ErrorHelper)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bootstrap.Initialize(app, db)

	app.Use(helpers.NotFoundHelper)

	app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
}
