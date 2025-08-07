package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
)

// AuthRoutes handles authentication-related routes using centralized dependencies
func AuthRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.GetAuthMiddleware()
	authHandler := dependencies.GetAuthHandler()

	// Auth routes
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	auth.Get("/profile", authMiddleware.JWTAuth(), authHandler.GetProfile)
}
