package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
)

// AuthRoutes handles authentication-related routes using centralized dependencies
func AuthRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.AuthMiddleware
	authHandler := dependencies.AuthHandler

	// Auth routes
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")

	// Apply stricter rate limiting for auth endpoints
	auth.Use(middleware.AuthRateLimiter())

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	// Protected routes
	auth.Get("/profile", authMiddleware.JWTAuth(), authHandler.GetProfile)
	auth.Put("/change-password", authMiddleware.JWTAuth(), authHandler.ChangePassword)
}
