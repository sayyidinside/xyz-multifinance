package shutdown

import (
	"context"
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
		timeout:     15 * time.Second,
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
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// Step 1: Shutdown HTTP server
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Step 2: Shutdown database
	s.shutdownDB(ctx)

	// Step 3: Close Redis clients
	if err := s.redisClient.CacheClient.Shutdown(); err != nil {
		log.Printf("Redis cache shutdown error: %v", err)
	}

	if err := s.redisClient.LockClient.Shutdown(); err != nil {
		log.Printf("Redis lock shutdown error: %v", err)
	}

	// Final status
	duration := time.Since(startTime)
	select {
	case <-ctx.Done():
		log.Printf("Shutdown timed out after %s", duration)
	default:
		log.Printf("Shutdown completed in %s", duration)
	}
}

func (s *ShutdownHandler) shutdownDB(ctx context.Context) {
	sqlDB, err := s.db.DB()
	if err != nil {
		log.Printf("Failed to get database instance: %v", err)
		return
	}

	// Prevent new connections
	sqlDB.SetMaxOpenConns(0)
	log.Println("Database: Stopped accepting new connections")

	// Check in-flight connections
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Database shutdown timeout - forcing close")
			if err := sqlDB.Close(); err != nil {
				log.Printf("Database force close error: %v", err)
			}
			return
		case <-ticker.C:
			stats := sqlDB.Stats()
			if stats.InUse > 0 {
				log.Printf("Database: Waiting on %d connections", stats.InUse)
			}
			if stats.InUse == 0 {
				if err := sqlDB.Close(); err != nil {
					log.Printf("Database close error: %v", err)
				} else {
					log.Println("Database: Closed successfully")
				}
				return
			}
		}
	}
}
