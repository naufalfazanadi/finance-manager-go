package container

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/handlers"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
	"gorm.io/gorm"
)

// ServiceContainer holds all application dependencies
type ServiceContainer struct {
	// Core dependencies
	DB        *gorm.DB
	Validator *validator.Validator

	// Repositories
	UserRepo   repositories.UserRepository
	WalletRepo repositories.WalletRepository

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Use cases
	AuthUseCase   usecases.AuthUseCaseInterface
	UserUseCase   usecases.UserUseCaseInterface
	WalletUseCase usecases.WalletUseCaseInterface

	// Handlers
	AuthHandler   *handlers.AuthHandler
	UserHandler   *handlers.UserHandler
	WalletHandler *handlers.WalletHandler
}

// NewServiceContainer creates and initializes all application dependencies
func NewServiceContainer(db *gorm.DB, validator *validator.Validator) *ServiceContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(userRepo)
	userUseCase := usecases.NewUserUseCase(userRepo)
	walletUseCase := usecases.NewWalletUseCase(walletRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase, validator)
	userHandler := handlers.NewUserHandler(userUseCase, validator)
	walletHandler := handlers.NewWalletHandler(walletUseCase, validator)

	return &ServiceContainer{
		DB:             db,
		Validator:      validator,
		UserRepo:       userRepo,
		WalletRepo:     walletRepo,
		AuthMiddleware: authMiddleware,
		AuthUseCase:    authUseCase,
		UserUseCase:    userUseCase,
		WalletUseCase:  walletUseCase,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		WalletHandler:  walletHandler,
	}
}

// GetAuthHandler returns the auth handler
func (sc *ServiceContainer) GetAuthHandler() *handlers.AuthHandler {
	return sc.AuthHandler
}

// GetUserHandler returns the user handler
func (sc *ServiceContainer) GetUserHandler() *handlers.UserHandler {
	return sc.UserHandler
}

// GetWalletHandler returns the wallet handler
func (sc *ServiceContainer) GetWalletHandler() *handlers.WalletHandler {
	return sc.WalletHandler
}

// GetAuthMiddleware returns the auth middleware
func (sc *ServiceContainer) GetAuthMiddleware() *middleware.AuthMiddleware {
	return sc.AuthMiddleware
}

// GetAuthUseCase returns the auth use case
func (sc *ServiceContainer) GetAuthUseCase() usecases.AuthUseCaseInterface {
	return sc.AuthUseCase
}

// GetUserUseCase returns the user use case
func (sc *ServiceContainer) GetUserUseCase() usecases.UserUseCaseInterface {
	return sc.UserUseCase
}

// GetWalletUseCase returns the wallet use case
func (sc *ServiceContainer) GetWalletUseCase() usecases.WalletUseCaseInterface {
	return sc.WalletUseCase
}

// GetDB returns the database instance
func (sc *ServiceContainer) GetDB() *gorm.DB {
	return sc.DB
}
