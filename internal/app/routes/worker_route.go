package routes

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"

	"github.com/gofiber/fiber/v2"
)

// WorkerRoutes handles worker-related routes using centralized dependencies
func WorkerRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.AuthMiddleware
	workerHandler := dependencies.WorkerHandler

	// Worker routes (admin only)
	v1 := api.Group("/v1")
	workers := v1.Group("/workers")

	// All worker routes require authentication
	workers.Get("/status", authMiddleware.JWTAuth(), workerHandler.GetWorkerStatus)           // Get worker status
	workers.Post("/balance-sync", authMiddleware.JWTAuth(), workerHandler.TriggerBalanceSync) // Trigger manual balance sync
}
