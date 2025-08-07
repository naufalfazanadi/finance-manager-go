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
	UserRepo repositories.UserRepository

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Use cases
	UserUseCase usecases.UserUseCaseInterface
	AuthUseCase usecases.AuthUseCaseInterface

	// Handlers
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
}

// NewServiceContainer creates and initializes all application dependencies
func NewServiceContainer(db *gorm.DB, validator *validator.Validator) *ServiceContainer {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo)
	authUseCase := usecases.NewAuthUseCase(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase, validator)
	authHandler := handlers.NewAuthHandler(authUseCase, validator)

	return &ServiceContainer{
		DB:             db,
		Validator:      validator,
		UserRepo:       userRepo,
		AuthMiddleware: authMiddleware,
		UserUseCase:    userUseCase,
		AuthUseCase:    authUseCase,
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
	}
}

// GetUserHandler returns the user handler
func (sc *ServiceContainer) GetUserHandler() *handlers.UserHandler {
	return sc.UserHandler
}

// GetAuthHandler returns the auth handler
func (sc *ServiceContainer) GetAuthHandler() *handlers.AuthHandler {
	return sc.AuthHandler
}

// GetAuthMiddleware returns the auth middleware
func (sc *ServiceContainer) GetAuthMiddleware() *middleware.AuthMiddleware {
	return sc.AuthMiddleware
}

// GetUserUseCase returns the user use case
func (sc *ServiceContainer) GetUserUseCase() usecases.UserUseCaseInterface {
	return sc.UserUseCase
}

// GetAuthUseCase returns the auth use case
func (sc *ServiceContainer) GetAuthUseCase() usecases.AuthUseCaseInterface {
	return sc.AuthUseCase
}

// GetDB returns the database instance
func (sc *ServiceContainer) GetDB() *gorm.DB {
	return sc.DB
}
