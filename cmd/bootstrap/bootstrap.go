package bootstrap

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/worker"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/redis"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

func Initialize(app *fiber.App, db *gorm.DB, cacheRedis *redis.CacheClient, lockRedis *redis.LockClient) {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	moduleRepo := repository.NewModuleRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	limitRepo := repository.NewLimitRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	installmentRepo := repository.NewIntallmentRepository(db)

	// Service
	userService := service.NewUserService(userRepo, roleRepo)
	permissionService := service.NewPermissionService(permissionRepo, moduleRepo)
	moduleService := service.NewModuleService(moduleRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	authService := service.NewAuthService(refreshTokenRepo, userRepo)
	registrationService := service.NewRegistrationService(userRepo, roleRepo, limitRepo)
	profileService := service.NewProfileService(userRepo, profileRepo)
	limitService := service.NewLimitService(userRepo, limitRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo, limitRepo, installmentRepo, lockRedis)

	// Handler
	userHandler := handler.NewUserHandler(userService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	moduleHandler := handler.NewModuleHandler(moduleService)
	roleHandler := handler.NewRoleHandler(roleService)
	authHandler := handler.NewAuthHandler(authService)
	registrationHandler := handler.NewRegistrationHandler(registrationService)
	profileHandler := handler.NewProfileHandler(profileService)
	limitHandler := handler.NewLimitHandler(limitService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Setup handler to send to routes setup
	handler := &handler.Handlers{
		UserManagementHandler: &handler.UserManagementHandler{
			UserHandler:       userHandler,
			PermissionHandler: permissionHandler,
			ModuleHandler:     moduleHandler,
			RoleHandler:       roleHandler,
			ProfileHandler:    profileHandler,
		},
		TransactionManagementHandler: &handler.TransactionManagementHandler{
			LimitHandler:       limitHandler,
			TransactionHandler: transactionHandler,
		},
		AuthHandler:         authHandler,
		RegistrationHandler: registrationHandler,
	}

	routes.Setup(app, handler)
}

func InitApp() {
	if err := config.LoadConfig(); err != nil {
		log.Println(err.Error())
	}

	worker.StartLogWorker()

	helpers.InitLogger()

	middleware.InitWhitelistIP()

}
