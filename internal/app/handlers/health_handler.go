package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/database"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/minio"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db          *gorm.DB
	minioClient minio.Client
}

func NewHealthHandler(db *gorm.DB, minioClient minio.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		minioClient: minioClient,
	}
}

// CheckHealth performs overall health check
// @Summary Basic health check
// @Description Basic health check endpoint to verify service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{"status": "ok", "service": "finance-manager-go", "message": "Service is running"}"
// @Router / [get]
// @Router /health [get]
func (h *HealthHandler) CheckHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "finance-manager-go",
		"message": "Service is running",
	})
}

// CheckDatabase performs database health check with connection stats
// @Summary Database health check
// @Description Database health check with connection statistics
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{"status": "ok", "database": "connected", "connection_stats": {...}}"
// @Failure 503 {object} map[string]interface{} "{"status": "error", "message": "Database connection failed", "error": "..."}"
// @Router /health/db [get]
func (h *HealthHandler) CheckDatabase(c *fiber.Ctx) error {
	// Direct health check without goroutine for debugging
	if err := database.HealthCheck(h.db); err != nil {
		logger.LogError(
			"HealthHandler.CheckDatabase",
			"Database health check failed",
			err,
		)
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
	}

	// Get connection stats
	stats, statsErr := database.GetConnectionStats(h.db)
	if statsErr != nil {
		logger.LogError(
			"HealthHandler.CheckDatabase",
			"Database connection stats unavailable but connection is healthy",
			statsErr,
		)
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"message":  "Database is healthy but stats unavailable",
		})
	}

	logger.LogSuccess(
		"HealthHandler.CheckDatabase",
		"Database health check successful with connection stats",
	)

	return c.JSON(fiber.Map{
		"status":   "ok",
		"database": "connected",
		"connection_stats": fiber.Map{
			"max_open_connections": stats.MaxOpenConnections,
			"open_connections":     stats.OpenConnections,
			"in_use":               stats.InUse,
			"idle":                 stats.Idle,
			"wait_count":           stats.WaitCount,
			"wait_duration":        stats.WaitDuration.String(),
			"max_idle_closed":      stats.MaxIdleClosed,
			"max_idle_time_closed": stats.MaxIdleTimeClosed,
			"max_lifetime_closed":  stats.MaxLifetimeClosed,
		},
	})
}

// CheckMinio performs Minio health check
// @Summary Minio health check
// @Description Minio health check to verify storage service connectivity
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{"status": "ok", "minio": "connected", "message": "Minio is healthy"}"
// @Failure 503 {object} map[string]interface{} "{"status": "error", "message": "Minio connection failed", "error": "..."}"
// @Router /health/minio [get]
func (h *HealthHandler) CheckMinio(c *fiber.Ctx) error {
	if h.minioClient == nil {
		logger.LogError(
			"HealthHandler.CheckMinio",
			"Minio health check failed - client not initialized",
			nil,
		)
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Minio client not initialized",
		})
	}

	// Create a context with timeout for the health check
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to perform a simple operation to check connectivity
	// We'll try to check if we can list buckets (this requires minimal permissions)
	testObject := minio.DownloadObject{
		BucketName:      "health-check", // This doesn't need to exist, just testing connectivity
		ObjectName:      "test",
		FilePath:        "",
		ExpiredInMinute: 1,
	}

	// This will fail but should give us connectivity status
	_, err := h.minioClient.GetObject(ctx, testObject)

	// If error is about bucket not existing, that means connection is working
	// If error is about connection, that means Minio is down
	if err != nil {
		// Check if it's a connection error or just bucket/object not found
		errStr := err.Error()
		if containsConnectionError(errStr) {
			logger.LogError(
				"HealthHandler.CheckMinio",
				"Minio connection failed during health check",
				err,
			)
			return c.Status(503).JSON(fiber.Map{
				"status":  "error",
				"message": "Minio connection failed",
				"error":   err.Error(),
			})
		}
		// Log that connection is working but bucket/object not found (expected)
		logger.LogSuccess(
			"HealthHandler.CheckMinio",
			"Minio connection successful - bucket/object not found as expected",
		)
	} else {
		// Unexpected success - log it
		logger.LogSuccess(
			"HealthHandler.CheckMinio",
			"Minio health check successful - unexpected object found",
		)
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"minio":   "connected",
		"message": "Minio is healthy",
	})
}

// containsConnectionError checks if the error is related to connection issues
func containsConnectionError(errStr string) bool {
	connectionErrors := []string{
		"connection refused",
		"connection timeout",
		"no such host",
		"network is unreachable",
		"context deadline exceeded",
	}

	for _, connErr := range connectionErrors {
		if len(errStr) > 0 && len(connErr) > 0 {
			// Simple substring check
			for i := 0; i <= len(errStr)-len(connErr); i++ {
				match := true
				for j := 0; j < len(connErr); j++ {
					if errStr[i+j] != connErr[j] {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}

// CheckAll performs comprehensive health check for all services
// @Summary Comprehensive health check
// @Description Health check for all services including database and Minio
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{"status": "ok", "services": {...}}"
// @Failure 503 {object} map[string]interface{} "{"status": "degraded", "services": {...}}"
// @Router /health/all [get]
func (h *HealthHandler) CheckAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	services := make(map[string]interface{})
	overallStatus := "ok"

	// Check database
	dbStatus := "ok"
	dbDetails := make(map[string]interface{})

	if err := database.HealthCheck(h.db); err != nil {
		dbStatus = "error"
		dbDetails["error"] = err.Error()
		overallStatus = "degraded"
	} else {
		stats, statsErr := database.GetConnectionStats(h.db)
		if statsErr == nil {
			dbDetails["connection_stats"] = fiber.Map{
				"max_open_connections": stats.MaxOpenConnections,
				"open_connections":     stats.OpenConnections,
				"in_use":               stats.InUse,
				"idle":                 stats.Idle,
			}
		}
	}

	services["database"] = fiber.Map{
		"status":  dbStatus,
		"details": dbDetails,
	}

	// Check Minio
	minioStatus := "ok"
	minioDetails := make(map[string]interface{})

	if h.minioClient == nil {
		minioStatus = "error"
		minioDetails["error"] = "Minio client not initialized"
		overallStatus = "degraded"
	} else {
		testObject := minio.DownloadObject{
			BucketName:      "health-check",
			ObjectName:      "test",
			FilePath:        "",
			ExpiredInMinute: 1,
		}

		_, err := h.minioClient.GetObject(ctx, testObject)
		if err != nil {
			errStr := err.Error()
			if containsConnectionError(errStr) {
				minioStatus = "error"
				minioDetails["error"] = err.Error()
				overallStatus = "degraded"
			} else {
				// Connection is working, just bucket/object not found
				minioDetails["message"] = "Connection successful"
			}
		}
	}

	services["minio"] = fiber.Map{
		"status":  minioStatus,
		"details": minioDetails,
	}

	statusCode := 200
	if overallStatus == "degraded" {
		statusCode = 503
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":    overallStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"services":  services,
	})
}
