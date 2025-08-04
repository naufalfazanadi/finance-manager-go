package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/naufalfazanadi/finance-manager-go/internal/app/routes"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/database"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.Init(cfg.App.LogLevel); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize database
	db := database.NewPostgresDB(cfg.Database)

	// Initialize validator
	validator := validator.New()

	// Setup routes with dependencies
	app := routes.Setup(db, validator)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("Starting server on " + serverAddr)

	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			logger.Error("Error during shutdown: " + err.Error())
		}
	}()
}
