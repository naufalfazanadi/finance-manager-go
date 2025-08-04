package routes

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, validator *validator.Validator) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())
	app.Use(middleware.CORS())

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Setup user routes with dependencies initialized per module
	UserRoutes(api, db, validator)

	return app
}
