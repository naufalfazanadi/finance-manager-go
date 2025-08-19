package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
)

// WalletRoutes handles wallet-related routes using centralized dependencies
func WalletRoutes(api fiber.Router, dependencies *container.ServiceContainer) {
	// Get handlers and middleware from centralized container
	authMiddleware := dependencies.AuthMiddleware
	walletHandler := dependencies.WalletHandler

	// Wallet routes
	v1 := api.Group("/v1")
	wallets := v1.Group("/wallets")

	// Protected routes (authentication required)
	wallets.Post("/", authMiddleware.JWTAuth(), walletHandler.CreateWallet)      // Create wallet (signup) - supports both JSON and multipart with optional photo
	wallets.Get("/", authMiddleware.JWTAuth(), walletHandler.GetWallets)         // Get all wallets (wallet/admin)
	wallets.Get("/:id", authMiddleware.JWTAuth(), walletHandler.GetWallet)       // Get wallet by ID (wallet/admin)
	wallets.Put("/:id", authMiddleware.JWTAuth(), walletHandler.UpdateWallet)    // Update wallet (wallet/admin) - supports both JSON and multipart with optional photo
	wallets.Delete("/:id", authMiddleware.JWTAuth(), walletHandler.DeleteWallet) // Soft delete wallet (admin only)
}
