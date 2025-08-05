package container

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/handlers"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
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

// NewContainer creates and initializes all application dependencies
func NewContainer(db *gorm.DB, validator *validator.Validator) *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(userRepo, validator)
	authUseCase := usecases.NewAuthUseCase(userRepo, validator)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userUseCase)
	authHandler := handlers.NewAuthHandler(authUseCase)

	return &Container{
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
func (c *Container) GetUserHandler() *handlers.UserHandler {
	return c.UserHandler
}

// GetAuthHandler returns the auth handler
func (c *Container) GetAuthHandler() *handlers.AuthHandler {
	return c.AuthHandler
}

// GetAuthMiddleware returns the auth middleware
func (c *Container) GetAuthMiddleware() *middleware.AuthMiddleware {
	return c.AuthMiddleware
}

// GetUserUseCase returns the user use case
func (c *Container) GetUserUseCase() usecases.UserUseCaseInterface {
	return c.UserUseCase
}

// GetAuthUseCase returns the auth use case
func (c *Container) GetAuthUseCase() usecases.AuthUseCaseInterface {
	return c.AuthUseCase
}
