package worker

import (
	"context"
	"time"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CronWorker struct {
	cron          *cron.Cron
	balanceSyncUC usecases.BalanceSyncUseCaseInterface
	db            *gorm.DB
	isRunning     bool
}

func NewCronWorker(balanceSyncUC usecases.BalanceSyncUseCaseInterface, db *gorm.DB) *CronWorker {
	// Create cron with logger and timezone
	c := cron.New(
		cron.WithLogger(cron.VerbosePrintfLogger(logger.Logger)),
		cron.WithLocation(time.UTC),
		cron.WithChain(cron.Recover(cron.VerbosePrintfLogger(logger.Logger))),
	)

	return &CronWorker{
		cron:          c,
		balanceSyncUC: balanceSyncUC,
		db:            db,
		isRunning:     false,
	}
}

// Start starts the cron worker and schedules jobs
func (w *CronWorker) Start() error {
	funcCtx := "CronWorker.Start"

	if w.isRunning {
		logger.LogSuccess(funcCtx, "Cron worker is already running", logrus.Fields{})
		return nil
	}

	// Schedule balance sync job to run every day at 00:00 AM UTC
	// Cron expression: "0 0 * * *" means minute=0, hour=0, every day, every month, every day of week
	_, err := w.cron.AddFunc("0 0 * * *", w.syncWalletBalances)
	if err != nil {
		logger.LogError(funcCtx, "failed to add balance sync cron job", err, logrus.Fields{})
		return err
	}

	// Optional: Add a test job that runs every minute for debugging (comment out in production)
	// _, err = w.cron.AddFunc("* * * * *", w.syncWalletBalances)
	// if err != nil {
	// 	logger.LogError(funcCtx, "failed to add test balance sync cron job", err, logrus.Fields{})
	// 	return err
	// }

	// Start the cron scheduler
	w.cron.Start()
	w.isRunning = true

	logger.LogSuccess(funcCtx, "Cron worker started successfully", logrus.Fields{
		"scheduled_jobs": len(w.cron.Entries()),
		"next_runs":      w.getNextRunTimes(),
	})

	return nil
}

// Stop stops the cron worker
func (w *CronWorker) Stop() {
	funcCtx := "CronWorker.Stop"

	if !w.isRunning {
		logger.LogSuccess(funcCtx, "Cron worker is not running", logrus.Fields{})
		return
	}

	// Stop the cron scheduler
	ctx := w.cron.Stop()
	w.isRunning = false

	// Wait for all running jobs to complete (with timeout)
	select {
	case <-ctx.Done():
		logger.LogSuccess(funcCtx, "All cron jobs completed gracefully", logrus.Fields{})
	case <-time.After(30 * time.Second):
		logger.LogSuccess(funcCtx, "Cron worker stopped (timeout waiting for jobs)", logrus.Fields{})
	}
}

// IsRunning returns whether the cron worker is running
func (w *CronWorker) IsRunning() bool {
	return w.isRunning
}

// GetStatus returns the status of the cron worker
func (w *CronWorker) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"is_running":     w.isRunning,
		"scheduled_jobs": len(w.cron.Entries()),
	}

	if w.isRunning {
		status["next_runs"] = w.getNextRunTimes()
	}

	return status
}

// syncWalletBalances is the job function that gets executed by cron
func (w *CronWorker) syncWalletBalances() {
	funcCtx := "CronWorker.syncWalletBalances"
	jobStart := time.Now()

	logger.LogSuccess(funcCtx, "Starting scheduled wallet balance sync job", logrus.Fields{
		"scheduled_time": jobStart.Format(time.RFC3339),
	})

	// Create context with timeout for the sync job
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Execute the balance sync using usecase
	err := w.balanceSyncUC.SyncAllWalletBalances(ctx)

	jobDuration := time.Since(jobStart)

	if err != nil {
		logger.LogError(funcCtx, "Scheduled wallet balance sync job failed", err, logrus.Fields{
			"job_duration":   jobDuration.String(),
			"scheduled_time": jobStart.Format(time.RFC3339),
		})
		return
	}

	logger.LogSuccess(funcCtx, "Scheduled wallet balance sync job completed successfully", logrus.Fields{
		"job_duration":   jobDuration.String(),
		"scheduled_time": jobStart.Format(time.RFC3339),
	})
}

// TriggerSync manually triggers balance sync for all wallets
func (w *CronWorker) TriggerSync(ctx context.Context) error {
	funcCtx := "CronWorker.TriggerSync"

	logger.LogSuccess(funcCtx, "Manual wallet balance sync triggered", logrus.Fields{})

	// Create context with timeout for manual sync
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	err := w.balanceSyncUC.SyncAllWalletBalances(timeoutCtx)
	if err != nil {
		logger.LogError(funcCtx, "Manual wallet balance sync failed", err, logrus.Fields{})
		return err
	}

	logger.LogSuccess(funcCtx, "Manual wallet balance sync completed successfully", logrus.Fields{})
	return nil
}

// getNextRunTimes returns the next scheduled run times for debugging
func (w *CronWorker) getNextRunTimes() []string {
	if !w.isRunning {
		return []string{}
	}

	var nextRuns []string
	for _, entry := range w.cron.Entries() {
		nextRuns = append(nextRuns, entry.Next.Format(time.RFC3339))
	}

	return nextRuns
}
