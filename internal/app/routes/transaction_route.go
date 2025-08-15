package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
)

// TransactionRoutes handles transaction-related routes using centralized dependencies
func TransactionRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.GetAuthMiddleware()
	transactionHandler := dependencies.GetTransactionHandler()

	// Transaction routes
	v1 := api.Group("/v1")
	transactions := v1.Group("/transactions")

	// Protected routes (authentication required)
	transactions.Post("/", authMiddleware.JWTAuth(), transactionHandler.CreateTransaction)      // Create transaction
	transactions.Get("/", authMiddleware.JWTAuth(), transactionHandler.GetTransactions)         // Get all transactions
	transactions.Get("/:id", authMiddleware.JWTAuth(), transactionHandler.GetTransaction)       // Get transaction by ID
	transactions.Put("/:id", authMiddleware.JWTAuth(), transactionHandler.UpdateTransaction)    // Update transaction
	transactions.Delete("/:id", authMiddleware.JWTAuth(), transactionHandler.DeleteTransaction) // Soft delete transaction
}
