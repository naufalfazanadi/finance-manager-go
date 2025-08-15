package routes

import (
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/handlers"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Setup(dependencies *container.ServiceContainer) *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(middleware.CORS())

	// Rate limiter middleware - Global rate limiting
	app.Use(middleware.DefaultRateLimiter())

	// Swagger endpoint
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check endpoints with database and Minio monitoring
	healthHandler := handlers.NewHealthHandler(dependencies.GetDB(), dependencies.GetMinioClient())

	// Basic health check
	app.Get("/", healthHandler.CheckHealth)
	app.Get("/health", healthHandler.CheckHealth)

	// Database health check with connection stats
	app.Get("/health/db", healthHandler.CheckDatabase)

	// Minio health check
	app.Get("/health/minio", healthHandler.CheckMinio)

	// Comprehensive health check for all services
	app.Get("/health/all", healthHandler.CheckAll)

	// API routes
	api := app.Group("/api")

	// API-specific rate limiter - More restrictive for API endpoints
	api.Use(middleware.APIRateLimiter())

	// Setup routes with centralized dependencies
	AuthRoutes(api, dependencies)
	UserRoutes(api, dependencies)
	WalletRoutes(api, dependencies)
	TransactionRoutes(api, dependencies)

	return app
}
