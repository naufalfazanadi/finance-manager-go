package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(dbConfig config.DatabaseConfig) *gorm.DB {
	var db *gorm.DB
	var err error

	log.Printf("Attempting to connect to database with pool settings: MaxOpen=%d, MaxIdle=%d",
		dbConfig.MaxOpenConns, dbConfig.MaxIdleConns)

	// Connection with retry mechanism
	for attempts := 1; attempts <= dbConfig.MaxRetries; attempts++ {
		db, err = connectWithPool(dbConfig)
		if err == nil {
			break
		}

		if attempts == dbConfig.MaxRetries {
			log.Fatalf("Failed to connect to database after %d attempts: %v", attempts, err)
		}

		log.Printf("Database connection attempt %d failed: %v. Retrying in %ds...",
			attempts, err, dbConfig.RetryDelay)
		time.Sleep(time.Duration(dbConfig.RetryDelay) * time.Second)
	}

	// Configure connection pool with goroutine-safe settings
	if err := configureConnectionPool(db, dbConfig); err != nil {
		log.Fatal("Failed to configure connection pool:", err)
	}

	// Test connection with timeout
	if err := testConnection(db, 10*time.Second); err != nil {
		log.Fatal("Failed to test database connection:", err)
	}

	// Auto migrate with timeout context
	if err := migrateWithTimeout(db, 30*time.Second); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Printf("Database connected successfully with connection pool (MaxOpen: %d, MaxIdle: %d, MaxLifetime: %dm)",
		dbConfig.MaxOpenConns, dbConfig.MaxIdleConns, dbConfig.ConnMaxLifetime)
	return db
}

func connectWithPool(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC connect_timeout=%d",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
		dbConfig.ConnectTimeout,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		// Optimizations for concurrent access
		PrepareStmt:                              true, // Prepare statements for better performance
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func configureConnectionPool(db *gorm.DB, dbConfig config.DatabaseConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// Connection Pool Configuration for Goroutines
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)                                    // Maximum open connections
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)                                    // Maximum idle connections
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Minute) // Connection lifetime
	sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTime) * time.Minute) // Idle timeout

	log.Printf("Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%dm, MaxIdleTime=%dm",
		dbConfig.MaxOpenConns, dbConfig.MaxIdleConns, dbConfig.ConnMaxLifetime, dbConfig.ConnMaxIdleTime)

	return nil
}

func testConnection(db *gorm.DB, timeout time.Duration) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Test with ping using goroutine for timeout
	pingChan := make(chan error, 1)
	go func() {
		pingChan <- sqlDB.PingContext(ctx)
	}()

	select {
	case err := <-pingChan:
		if err != nil {
			return fmt.Errorf("database ping failed: %w", err)
		}
		log.Println("Database connection test successful")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("database ping timeout after %v", timeout)
	}
}

func migrateWithTimeout(db *gorm.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Run migration in goroutine with timeout
	migrationChan := make(chan error, 1)
	go func() {
		err := db.AutoMigrate(
			&entities.User{},
			&entities.Wallet{},
			// Add other entities here as your project grows
		)
		migrationChan <- err
	}()

	select {
	case err := <-migrationChan:
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		log.Println("Database migration completed successfully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("database migration timeout after %v", timeout)
	}
}

// HealthCheck performs a health check with goroutine safety
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// GetConnectionStats returns connection pool statistics
func GetConnectionStats(db *gorm.DB) (*sql.DBStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return &stats, nil
}
