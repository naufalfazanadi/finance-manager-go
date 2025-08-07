package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/database"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
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
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
	}

	// Get connection stats
	stats, statsErr := database.GetConnectionStats(h.db)
	if statsErr != nil {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
			"message":  "Database is healthy but stats unavailable",
		})
	}

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
