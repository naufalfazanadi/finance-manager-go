package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/worker"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

type WorkerHandler struct {
	cronWorker *worker.CronWorker
}

func NewWorkerHandler(cronWorker *worker.CronWorker) *WorkerHandler {
	return &WorkerHandler{
		cronWorker: cronWorker,
	}
}

// GetWorkerStatus gets the current status of the worker
// @Summary Get worker status
// @Description Get the current status of the balance sync worker
// @Tags Worker
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /api/v1/worker/status [get]
func (h *WorkerHandler) GetWorkerStatus(c *fiber.Ctx) error {
	funcCtx := "WorkerHandler.GetWorkerStatus"

	status := h.cronWorker.GetStatus()

	logger.LogSuccess(funcCtx, "Retrieved worker status", logrus.Fields{
		"is_running": status["is_running"],
	})

	return helpers.SuccessResponse(c, "Worker status retrieved successfully", status)
}

// TriggerBalanceSync manually triggers balance sync for all wallets
// @Summary Trigger balance sync
// @Description Manually trigger balance sync for all wallets
// @Tags Worker
// @Accept json
// @Produce json
// @Success 200 {object} helpers.Response
// @Router /api/v1/worker/balance-sync [post]
func (h *WorkerHandler) TriggerBalanceSync(c *fiber.Ctx) error {
	funcCtx := "WorkerHandler.TriggerBalanceSync"

	logger.LogSuccess(funcCtx, "Manual balance sync triggered", logrus.Fields{})

	err := h.cronWorker.TriggerSync(c.Context())
	if err != nil {
		logger.LogError(funcCtx, "Failed to trigger balance sync", err, logrus.Fields{})
		return helpers.InternalServerErrorResponse(c, "Failed to trigger balance sync", err.Error())
	}

	logger.LogSuccess(funcCtx, "Balance sync completed successfully", logrus.Fields{})
	return helpers.SuccessResponse(c, "Balance sync completed successfully", nil)
}
