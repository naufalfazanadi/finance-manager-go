package container

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/app/handlers"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/middleware"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/internal/worker"
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
	BalanceSyncUseCase usecases.BalanceSyncUseCaseInterface

	// Workers
	CronWorker *worker.CronWorker

	// Handlers
	AuthHandler        *handlers.AuthHandler
	UserHandler        *handlers.UserHandler
	WalletHandler      *handlers.WalletHandler
	TransactionHandler *handlers.TransactionHandler
	WorkerHandler      *handlers.WorkerHandler
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
	balanceSyncUseCase := usecases.NewBalanceSyncUseCase(walletRepo, transactionRepo, db)

	// Initialize workers
	cronWorker := worker.NewCronWorker(balanceSyncUseCase, db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase, validator)
	userHandler := handlers.NewUserHandler(userUseCase, validator)
	walletHandler := handlers.NewWalletHandler(walletUseCase, validator)
	transactionHandler := handlers.NewTransactionHandler(transactionUseCase, validator)
	workerHandler := handlers.NewWorkerHandler(cronWorker)

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
		BalanceSyncUseCase: balanceSyncUseCase,
		CronWorker:         cronWorker,
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		WalletHandler:      walletHandler,
		TransactionHandler: transactionHandler,
		WorkerHandler:      workerHandler,
	}
}
