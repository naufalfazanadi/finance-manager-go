package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
)

// DashboardRoutes handles dashboard-related routes using centralized dependencies
func DashboardRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.AuthMiddleware
	dashboardHandler := dependencies.DashboardHandler

	// Dashboard routes
	v1 := api.Group("/v1")
	dashboard := v1.Group("/dashboard")

	// Protected routes (authentication required)
	dashboard.Get("/users/:id/monthly-summary", authMiddleware.JWTAuth(), dashboardHandler.GetMonthlySumByUser) // Get monthly sum by user ID
}
