package usecases

import (
	"context"
	"fmt"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BalanceSyncUseCaseInterface interface {
	SyncAllWalletBalances(ctx context.Context) error
	SyncWalletBalance(ctx context.Context, wallet *entities.Wallet) error
}

type BalanceSyncUseCase struct {
	db              *gorm.DB
	walletRepo      repositories.WalletRepository
	transactionRepo repositories.TransactionRepository
}

func NewBalanceSyncUseCase(walletRepo repositories.WalletRepository, transactionRepo repositories.TransactionRepository, db *gorm.DB) BalanceSyncUseCaseInterface {
	return &BalanceSyncUseCase{
		db:              db,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

// SyncAllWalletBalances recalculates and updates all wallet balances based on active transactions
func (uc *BalanceSyncUseCase) SyncAllWalletBalances(ctx context.Context) error {
	funcCtx := "BalanceSyncUseCase.SyncAllWalletBalances"

	logger.LogSuccess(funcCtx, "Starting wallet balance sync for all wallets", logrus.Fields{})

	// Get all active wallets using the correct method signature
	queryParams := &dto.QueryParams{
		PaginationQuery: &dto.PaginationQuery{
			Page:  1,
			Limit: 1000, // Process in batches if needed
		},
		FilterQuery: &dto.FilterQuery{},
	}

	wallets, err := uc.walletRepo.GetAll(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to get wallets", err, logrus.Fields{})
		return fmt.Errorf("failed to get wallets: %w", err)
	}

	var syncErrors []error
	syncedCount := 0

	// Process each wallet
	for _, wallet := range wallets {
		if err := uc.SyncWalletBalance(ctx, wallet); err != nil {
			logger.LogError(funcCtx, "failed to sync wallet balance", err, logrus.Fields{
				"wallet_id":   wallet.ID.String(),
				"wallet_name": wallet.Name,
			})
			syncErrors = append(syncErrors, fmt.Errorf("wallet %s: %w", wallet.ID.String(), err))
			continue
		}
		syncedCount++
	}

	logger.LogSuccess(funcCtx, "Completed wallet balance sync", logrus.Fields{
		"total_wallets": len(wallets),
		"synced_count":  syncedCount,
		"error_count":   len(syncErrors),
	})

	if len(syncErrors) > 0 {
		return fmt.Errorf("sync completed with %d errors: %v", len(syncErrors), syncErrors)
	}

	return nil
}

// SyncWalletBalance performs the actual balance sync for a wallet
func (uc *BalanceSyncUseCase) SyncWalletBalance(ctx context.Context, wallet *entities.Wallet) error {
	funcCtx := "BalanceSyncUseCase.SyncWalletBalance"

	// Get all active transactions for this wallet using the correct method
	transactions, err := uc.transactionRepo.GetByWalletID(ctx, wallet.ID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get transactions for wallet", err, logrus.Fields{
			"wallet_id": wallet.ID.String(),
		})
		return fmt.Errorf("failed to get transactions for wallet: %w", err)
	}

	// Calculate the correct balance based on transactions
	calculatedBalance := 0.0
	for _, transaction := range transactions {
		if !transaction.IsDeleted && !transaction.DeletedAt.Valid {
			calculatedBalance += transaction.GetWalletImpact()
		}
	}

	// Check if balance needs updating
	if wallet.Balance == calculatedBalance {
		logger.LogSuccess(funcCtx, "Wallet balance already correct", logrus.Fields{
			"wallet_id":          wallet.ID.String(),
			"wallet_name":        wallet.Name,
			"balance":            wallet.Balance,
			"calculated_balance": calculatedBalance,
		})
		return nil
	}

	// Update wallet balance
	oldBalance := wallet.Balance
	wallet.Balance = calculatedBalance

	if err := uc.walletRepo.Update(ctx, wallet); err != nil {
		logger.LogError(funcCtx, "failed to update wallet balance", err, logrus.Fields{
			"wallet_id":          wallet.ID.String(),
			"old_balance":        oldBalance,
			"calculated_balance": calculatedBalance,
		})
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	logger.LogSuccess(funcCtx, "Successfully synced wallet balance", logrus.Fields{
		"wallet_id":         wallet.ID.String(),
		"wallet_name":       wallet.Name,
		"old_balance":       oldBalance,
		"new_balance":       calculatedBalance,
		"difference":        calculatedBalance - oldBalance,
		"transaction_count": len(transactions),
	})

	return nil
}
