package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionUseCaseInterface interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetTransaction(ctx context.Context, id uuid.UUID, loggedUserID uuid.UUID) (*dto.TransactionResponse, error)
	GetTransactions(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.TransactionResponse], error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionResponse, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error // This now does soft delete
}

type TransactionUseCase struct {
	transactionRepo repositories.TransactionRepository
	walletRepo      repositories.WalletRepository
	userRepo        repositories.UserRepository
	db              *gorm.DB
}

func NewTransactionUseCase(
	transactionRepo repositories.TransactionRepository,
	walletRepo repositories.WalletRepository,
	userRepo repositories.UserRepository,
	db *gorm.DB,
) TransactionUseCaseInterface {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		userRepo:        userRepo,
		db:              db,
	}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	funcCtx := "CreateTransaction"

	// Start transaction to ensure consistency
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify user exists
	_, err := uc.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "user not found", err, logrus.Fields{"user_id": req.UserID.String()})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	// Verify wallet exists and belongs to user (or user is admin)
	wallet, err := uc.walletRepo.GetByID(ctx, req.WalletID)
	if err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "wallet not found", err, logrus.Fields{"wallet_id": req.WalletID.String()})
		return nil, helpers.NewNotFoundError("wallet not found", "")
	}

	// Check if user owns the wallet
	if wallet.UserID != req.UserID {
		tx.Rollback()
		logger.LogError(funcCtx, "wallet does not belong to user", nil, logrus.Fields{
			"wallet_id":       req.WalletID.String(),
			"wallet_user_id":  wallet.UserID.String(),
			"request_user_id": req.UserID.String(),
		})
		return nil, helpers.NewForbiddenError("wallet does not belong to the specified user", "")
	}

	// Create transaction entity
	transaction := &entities.Transaction{
		Name:      req.Name,
		Cost:      req.Cost,
		Type:      entities.TransactionType(req.Type),
		Note:      req.Note,
		TCategory: req.TCategory,
		UserID:    req.UserID,
		WalletID:  req.WalletID,
	}

	// Save transaction within transaction
	transactionRepo := repositories.NewTransactionRepository(tx)
	if err := transactionRepo.Create(ctx, transaction); err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to create transaction", err, logrus.Fields{
			"name":      req.Name,
			"user_id":   req.UserID.String(),
			"wallet_id": req.WalletID.String(),
		})
		return nil, helpers.NewInternalError("failed to create transaction", err.Error())
	}

	// Update wallet balance using the GetWalletImpact method
	wallet.Balance += transaction.GetWalletImpact()
	walletRepo := repositories.NewWalletRepository(tx)
	if err := walletRepo.Update(ctx, wallet); err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to update wallet balance", err, logrus.Fields{
			"wallet_id":     req.WalletID.String(),
			"wallet_impact": transaction.GetWalletImpact(),
		})
		return nil, helpers.NewInternalError("failed to update wallet balance", err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.LogError(funcCtx, "failed to commit transaction", err, logrus.Fields{
			"transaction_id": transaction.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to commit transaction", err.Error())
	}

	// Reload transaction with relationships
	createdTransaction, err := uc.transactionRepo.GetByID(ctx, transaction.ID)
	if err != nil {
		logger.LogError(funcCtx, "failed to reload created transaction", err, logrus.Fields{
			"transaction_id": transaction.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to reload created transaction", err.Error())
	}

	return dto.MapToTransactionResponse(createdTransaction), nil
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, id uuid.UUID, loggedUserID uuid.UUID) (*dto.TransactionResponse, error) {
	funcCtx := "GetTransaction"

	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get transaction", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("transaction not found", "")
	}

	// Authorization check: non-admin users can only access their own transactions
	if loggedUserID != transaction.UserID {
		logger.LogError(funcCtx, "unauthorized access to transaction", nil, logrus.Fields{
			"transaction_id":      id.String(),
			"transaction_user_id": transaction.UserID.String(),
			"logged_user_id":      loggedUserID.String(),
		})
		return nil, helpers.NewNotFoundError("transaction not found", "")
	}

	return dto.MapToTransactionResponse(transaction), nil
}

func (uc *TransactionUseCase) GetTransactions(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.TransactionResponse], error) {
	funcCtx := "GetTransactions"

	transactions, err := uc.transactionRepo.GetAll(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to get transactions", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to get transactions", err.Error())
	}

	total, err := uc.transactionRepo.CountWithFilters(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to count transactions", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to count transactions", err.Error())
	}

	transactionResponses := make([]dto.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionResponses[i] = *dto.MapToTransactionResponse(transaction)
	}

	paginationMeta := helpers.NewPaginationMeta(queryParams.Page, queryParams.Limit, total)

	return &dto.PaginationData[dto.TransactionResponse]{
		Data: transactionResponses,
		Meta: paginationMeta,
	}, nil
}

func (uc *TransactionUseCase) UpdateTransaction(ctx context.Context, id uuid.UUID, req *dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	funcCtx := "UpdateTransaction"

	// Start transaction to ensure consistency
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get existing transaction
	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to get transaction", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("transaction not found", "")
	}

	// Store original values for balance calculation
	originalCost := transaction.Cost
	originalType := transaction.Type
	originalWalletID := transaction.WalletID

	// Update transaction fields
	if req.Name != "" {
		transaction.Name = req.Name
	}
	if req.Note != "" {
		transaction.Note = req.Note
	}
	if req.TCategory != "" {
		transaction.TCategory = req.TCategory
	}

	// Handle cost update
	costChanged := false
	if req.Cost > 0 && req.Cost != originalCost {
		transaction.Cost = req.Cost
		costChanged = true
	}

	// Handle type update
	typeChanged := false
	if req.Type != "" && req.Type != string(originalType) {
		transaction.Type = entities.TransactionType(req.Type)
		typeChanged = true
	}

	// Handle wallet change
	walletChanged := false
	if req.WalletID != uuid.Nil && req.WalletID != originalWalletID {
		// Verify new wallet exists and belongs to user
		newWallet, err := uc.walletRepo.GetByID(ctx, req.WalletID)
		if err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "new wallet not found", err, logrus.Fields{"wallet_id": req.WalletID.String()})
			return nil, helpers.NewNotFoundError("new wallet not found", "")
		}

		if newWallet.UserID != transaction.UserID {
			tx.Rollback()
			logger.LogError(funcCtx, "new wallet does not belong to user", nil, logrus.Fields{
				"wallet_id":           req.WalletID.String(),
				"wallet_user_id":      newWallet.UserID.String(),
				"transaction_user_id": transaction.UserID.String(),
			})
			return nil, helpers.NewForbiddenError("new wallet does not belong to the transaction user", "")
		}

		transaction.WalletID = req.WalletID
		walletChanged = true
	}

	// Handle user change (typically admin only)
	if req.UserID != uuid.Nil && req.UserID != transaction.UserID {
		// Verify user exists
		_, err := uc.userRepo.GetByID(ctx, req.UserID)
		if err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "new user not found", err, logrus.Fields{"user_id": req.UserID.String()})
			return nil, helpers.NewNotFoundError("new user not found", "")
		}
		transaction.UserID = req.UserID
	}

	// Update wallet balances if needed
	walletRepo := repositories.NewWalletRepository(tx)

	if walletChanged {
		// Wallet changed: reverse from old wallet, apply to new wallet

		// Calculate original wallet impact
		originalTransaction := &entities.Transaction{
			Cost: originalCost,
			Type: originalType,
		}

		// Get original wallet and reverse the original transaction
		originalWallet, err := uc.walletRepo.GetByID(ctx, originalWalletID)
		if err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "failed to get original wallet", err, logrus.Fields{"wallet_id": originalWalletID.String()})
			return nil, helpers.NewInternalError("failed to get original wallet", err.Error())
		}

		// Reverse original impact from original wallet
		originalWallet.Balance -= originalTransaction.GetWalletImpact()
		if err := walletRepo.Update(ctx, originalWallet); err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "failed to reverse balance from original wallet", err, logrus.Fields{
				"wallet_id":      originalWalletID.String(),
				"impact_reverse": -originalTransaction.GetWalletImpact(),
			})
			return nil, helpers.NewInternalError("failed to reverse balance from original wallet", err.Error())
		}

		// Apply new impact to new wallet
		newWallet, err := uc.walletRepo.GetByID(ctx, transaction.WalletID)
		if err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "failed to get new wallet", err, logrus.Fields{"wallet_id": transaction.WalletID.String()})
			return nil, helpers.NewInternalError("failed to get new wallet", err.Error())
		}

		newWallet.Balance += transaction.GetWalletImpact()
		if err := walletRepo.Update(ctx, newWallet); err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "failed to apply balance to new wallet", err, logrus.Fields{
				"wallet_id":    transaction.WalletID.String(),
				"impact_apply": transaction.GetWalletImpact(),
			})
			return nil, helpers.NewInternalError("failed to apply balance to new wallet", err.Error())
		}

	} else if costChanged || typeChanged {
		// Same wallet, but cost or type changed: calculate difference and apply
		originalTransaction := &entities.Transaction{
			Cost: originalCost,
			Type: originalType,
		}

		impactDifference := transaction.GetWalletImpact() - originalTransaction.GetWalletImpact()

		wallet, err := uc.walletRepo.GetByID(ctx, transaction.WalletID)
		if err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "failed to get wallet for update", err, logrus.Fields{"wallet_id": transaction.WalletID.String()})
			return nil, helpers.NewInternalError("failed to get wallet for update", err.Error())
		}

		wallet.Balance += impactDifference
		if err := walletRepo.Update(ctx, wallet); err != nil {
			tx.Rollback()
			logger.LogError(funcCtx, "failed to update wallet balance", err, logrus.Fields{
				"wallet_id":         transaction.WalletID.String(),
				"impact_difference": impactDifference,
			})
			return nil, helpers.NewInternalError("failed to update wallet balance", err.Error())
		}
	}

	// Save updated transaction
	transactionRepo := repositories.NewTransactionRepository(tx)
	if err := transactionRepo.Update(ctx, transaction); err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to update transaction", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return nil, helpers.NewInternalError("failed to update transaction", err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.LogError(funcCtx, "failed to commit transaction update", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return nil, helpers.NewInternalError("failed to commit transaction update", err.Error())
	}

	// Reload transaction with relationships
	updatedTransaction, err := uc.transactionRepo.GetByID(ctx, transaction.ID)
	if err != nil {
		logger.LogError(funcCtx, "failed to reload updated transaction", err, logrus.Fields{
			"transaction_id": transaction.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to reload updated transaction", err.Error())
	}

	return dto.MapToTransactionResponse(updatedTransaction), nil
}

func (uc *TransactionUseCase) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	funcCtx := "DeleteTransaction"

	// Start transaction to ensure consistency
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get transaction to reverse wallet balance
	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to get transaction", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return helpers.NewNotFoundError("transaction not found", "")
	}

	// Get wallet and reverse the transaction cost
	wallet, err := uc.walletRepo.GetByID(ctx, transaction.WalletID)
	if err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to get wallet", err, logrus.Fields{
			"wallet_id": transaction.WalletID.String(),
		})
		return helpers.NewInternalError("failed to get wallet", err.Error())
	}

	// Reverse the transaction impact from wallet balance
	wallet.Balance -= transaction.GetWalletImpact()
	walletRepo := repositories.NewWalletRepository(tx)
	if err := walletRepo.Update(ctx, wallet); err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to reverse wallet balance", err, logrus.Fields{
			"wallet_id":         transaction.WalletID.String(),
			"impact_to_reverse": transaction.GetWalletImpact(),
		})
		return helpers.NewInternalError("failed to reverse wallet balance", err.Error())
	}

	// Use soft delete for default delete operation
	transactionRepo := repositories.NewTransactionRepository(tx)
	if err := transactionRepo.SoftDelete(ctx, id); err != nil {
		tx.Rollback()
		logger.LogError(funcCtx, "failed to delete transaction", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return helpers.NewInternalError("failed to delete transaction", err.Error())
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.LogError(funcCtx, "failed to commit transaction deletion", err, logrus.Fields{
			"transaction_id": id.String(),
		})
		return helpers.NewInternalError("failed to commit transaction deletion", err.Error())
	}

	return nil
}
