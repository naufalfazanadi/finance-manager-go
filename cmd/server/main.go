// @title Finance Manager API
// @version 1.0
// @description This is a finance manager API server built with Go and Fiber framework.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email naufalfazanadi@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/naufalfazanadi/finance-manager-go/docs"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/container"
	"github.com/naufalfazanadi/finance-manager-go/internal/app/routes"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/database"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.Init(cfg.App.LogLevel); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize DataDog tracer
	appEnv := cfg.App.Env
	if appEnv == "staging" || appEnv == "production" {
		tracer.Start(
			tracer.WithEnv(appEnv),
			tracer.WithServiceName("finance-manager-service"),
		)

		defer tracer.Stop()
	}

	// Initialize database
	db := database.NewPostgresDB(cfg.Database)

	// Initialize validator
	validator := validator.New()

	// Initialize centralized dependency container
	dependencies := container.NewServiceContainer(db, validator)

	// Start background workers
	cronWorker := dependencies.CronWorker
	cronWorker.Start()
	logger.Info("Cron worker started successfully")

	// Setup routes with dependencies
	app := routes.Setup(dependencies)

	// Prepare server address
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server on " + serverAddr)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			logger.Fatal("Failed to start server: " + err.Error())
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info("Shutting down server and workers...")

	// Stop background workers first
	cronWorker.Stop()

	// Then stop the server
	if err := app.Shutdown(); err != nil {
		logger.Error("Error during server shutdown: " + err.Error())
	}

	logger.Info("Server and workers shut down successfully")
}
