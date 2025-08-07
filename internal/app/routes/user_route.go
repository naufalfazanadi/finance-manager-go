package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
)

// UserRoutes handles user-related routes using centralized dependencies
func UserRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.GetAuthMiddleware()
	userHandler := dependencies.GetUserHandler()

	// User routes
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	// Public routes (no authentication required)
	users.Post("/", userHandler.CreateUser) // Create user (signup)

	// Protected routes (authentication required)
	users.Get("/", authMiddleware.JWTAuth(), userHandler.GetUsers)                                    // Get all users (user/admin)
	users.Get("/:id", authMiddleware.JWTAuth(), userHandler.GetUser)                                  // Get user by ID (user/admin)
	users.Put("/:id", authMiddleware.JWTAuth(), userHandler.UpdateUser)                               // Update user (user/admin)
	users.Delete("/:id", authMiddleware.JWTAuth(), middleware.RequireAdmin(), userHandler.DeleteUser) // Soft delete user (admin only)

	// Soft delete management routes (admin only)
	users.Patch("/:id/restore", authMiddleware.JWTAuth(), middleware.RequireAdmin(), userHandler.RestoreUser)  // Restore soft deleted user
	users.Delete("/:id/hard", authMiddleware.JWTAuth(), middleware.RequireAdmin(), userHandler.HardDeleteUser) // Hard delete user permanently

	// Routes to view deleted users (admin only)
	users.Get("/with-deleted", authMiddleware.JWTAuth(), middleware.RequireAdmin(), userHandler.GetUsersWithDeleted) // Get all users including deleted
	users.Get("/deleted", authMiddleware.JWTAuth(), middleware.RequireAdmin(), userHandler.GetOnlyDeletedUsers)      // Get only deleted users
}
