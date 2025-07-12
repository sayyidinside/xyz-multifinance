package tests

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"gorm.io/gorm"
)

var (
	TestApp *fiber.App
	TestDB  *gorm.DB
)

func TestMain(m *testing.M) {
	// Setup test environment
	setupTestEnvironment()

	// Run tests
	exitCode := m.Run()

	// Teardown
	teardownTestEnvironment()

	os.Exit(exitCode)
}

func setupTestEnvironment() {
	bootstrap.InitApp()

	config.AppConfig.Env = "test"
	// Initialize test database
	var err error
	TestDB, err = database.Connect()
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	// Create test app
	TestApp = fiber.New()
	bootstrap.Initialize(TestApp, TestDB, nil, nil)
}

func teardownTestEnvironment() {
	if err := TestDB.Migrator().DropTable(
		&entity.Module{},
		&entity.Permission{},
		&entity.Role{},
		&entity.RolePermission{},
		&entity.User{},
		&entity.UserProfile{},
		&entity.UserDocument{},
		&entity.RefreshToken{},
		&entity.Limit{},
		&entity.Transaction{},
		&entity.TransactionInstallment{},
		&entity.Payment{},
	); err != nil {
		panic("failed to clean database: " + err.Error())
	}

	// Close database
	if sqlDB, err := TestDB.DB(); err == nil {
		sqlDB.Close()
	}
}
