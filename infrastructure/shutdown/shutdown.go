package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/redis"
	"gorm.io/gorm"
)

type ShutdownHandler struct {
	app         *fiber.App
	db          *gorm.DB
	redisClient *redis.RedisClient
	timeout     time.Duration
}

func NewShutdownHandler(
	app *fiber.App,
	db *gorm.DB,
	redisClient *redis.RedisClient,
) *ShutdownHandler {
	return &ShutdownHandler{
		app:         app,
		db:          db,
		redisClient: redisClient,
		timeout:     15 * time.Second, // Default timeout
	}
}

func (s *ShutdownHandler) WithTimeout(duration time.Duration) *ShutdownHandler {
	s.timeout = duration
	return s
}

func (s *ShutdownHandler) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down...")

	// Step 1: Shutdown HTTP server
	if err := s.app.Shutdown(); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Step 2: Close database connection
	if sqlDB, err := s.db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Database shutdown error: %v", err)
		}
	}

	// Step 3: Close Redis clients
	if err := s.redisClient.CacheClient.Shutdown(); err != nil {
		log.Printf("Redis cache shutdown error: %v", err)
	}

	if err := s.redisClient.LockClient.Shutdown(); err != nil {
		log.Printf("Redis lock shutdown error: %v", err)
	}

	log.Println("Shutdown complete")
	os.Exit(0)
}
