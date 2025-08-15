package container

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/handlers"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/minio"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
	"gorm.io/gorm"
)

// ServiceContainer holds all application dependencies
type ServiceContainer struct {
	// Core dependencies
	DB          *gorm.DB
	Validator   *validator.Validator
	MinioClient minio.Client

	// Repositories
	UserRepo        repositories.UserRepository
	WalletRepo      repositories.WalletRepository
	TransactionRepo repositories.TransactionRepository

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Use cases
	AuthUseCase        usecases.AuthUseCaseInterface
	UserUseCase        usecases.UserUseCaseInterface
	WalletUseCase      usecases.WalletUseCaseInterface
	TransactionUseCase usecases.TransactionUseCaseInterface

	// Handlers
	AuthHandler        *handlers.AuthHandler
	UserHandler        *handlers.UserHandler
	WalletHandler      *handlers.WalletHandler
	TransactionHandler *handlers.TransactionHandler
}

// NewServiceContainer creates and initializes all application dependencies
func NewServiceContainer(db *gorm.DB, validator *validator.Validator) *ServiceContainer {
	// Initialize MinIO client
	minioClient, err := minio.NewClient()
	if err != nil {
		// Log error with context and details
		logger.LogError(
			"ServiceContainer.NewServiceContainer",
			"Failed to initialize Minio client - application will continue without Minio functionality",
			err,
		)
		minioClient = nil
	} else {
		// Log successful Minio client initialization
		logger.LogSuccess(
			"ServiceContainer.NewServiceContainer",
			"Minio client initialized successfully",
		)
	} // Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(userRepo)
	userUseCase := usecases.NewUserUseCase(userRepo)
	walletUseCase := usecases.NewWalletUseCase(walletRepo, userRepo)
	transactionUseCase := usecases.NewTransactionUseCase(transactionRepo, walletRepo, userRepo, db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase, validator)
	userHandler := handlers.NewUserHandler(userUseCase, validator)
	walletHandler := handlers.NewWalletHandler(walletUseCase, validator)
	transactionHandler := handlers.NewTransactionHandler(transactionUseCase, validator)

	// Log successful service container initialization
	logger.LogSuccess(
		"ServiceContainer.NewServiceContainer",
		"Service container initialized successfully with all dependencies",
	)

	return &ServiceContainer{
		DB:                 db,
		Validator:          validator,
		MinioClient:        minioClient,
		UserRepo:           userRepo,
		WalletRepo:         walletRepo,
		TransactionRepo:    transactionRepo,
		AuthMiddleware:     authMiddleware,
		AuthUseCase:        authUseCase,
		UserUseCase:        userUseCase,
		WalletUseCase:      walletUseCase,
		TransactionUseCase: transactionUseCase,
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		WalletHandler:      walletHandler,
		TransactionHandler: transactionHandler,
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

// GetTransactionHandler returns the transaction handler
func (sc *ServiceContainer) GetTransactionHandler() *handlers.TransactionHandler {
	return sc.TransactionHandler
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

// GetTransactionUseCase returns the transaction use case
func (sc *ServiceContainer) GetTransactionUseCase() usecases.TransactionUseCaseInterface {
	return sc.TransactionUseCase
}

// GetDB returns the database instance
func (sc *ServiceContainer) GetDB() *gorm.DB {
	return sc.DB
}

// GetMinioClient returns the MinIO client instance
func (sc *ServiceContainer) GetMinioClient() minio.Client {
	return sc.MinioClient
}
