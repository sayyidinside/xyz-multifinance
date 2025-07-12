package database

import (
	"fmt"
	"log"

	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	var dbHost, dbPort, dbUser, dbPassword, dbName string

	cfg := config.AppConfig

	switch cfg.Env {
	case "production":
		dbHost = cfg.ProdDbHost
		dbPort = cfg.ProdDbPort
		dbUser = cfg.ProdDbUsername
		dbPassword = cfg.ProdDbPassword
		dbName = cfg.ProdDbName
	case "development":
		dbHost = cfg.DevDbHost
		dbPort = cfg.DevDbPort
		dbUser = cfg.DevDbUsername
		dbPassword = cfg.DevDbPassword
		dbName = cfg.DevDbName
	case "test":
		dbHost = cfg.TestDbHost
		dbPort = cfg.TestDbPort
		dbUser = cfg.TestDbUsername
		dbPassword = cfg.TestDbPassword
		dbName = cfg.TestDbName
	default:
		dbHost = cfg.LocalDbHost
		dbPort = cfg.LocalDbPort
		dbUser = cfg.LocalDbUsername
		dbPassword = cfg.LocalDbPassword
		dbName = cfg.LocalDbName
	}

	// Construct database URL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	// Migrate and seed the database
	Migrate(db)
	Seeding(db)

	return db, nil
}
