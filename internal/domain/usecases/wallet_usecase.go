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
)

type WalletUseCaseInterface interface {
	CreateWallet(ctx context.Context, req *dto.CreateWalletRequest) (*dto.WalletResponse, error)
	GetWallet(ctx context.Context, id uuid.UUID, loggedUserID uuid.UUID) (*dto.WalletResponse, error)
	GetWallets(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.WalletResponse], error)
	UpdateWallet(ctx context.Context, id uuid.UUID, req *dto.UpdateWalletRequest) (*dto.WalletResponse, error)
	DeleteWallet(ctx context.Context, id uuid.UUID) error // This now does soft delete
}

type WalletUseCase struct {
	walletRepo repositories.WalletRepository
	userRepo   repositories.UserRepository
}

func NewWalletUseCase(walletRepo repositories.WalletRepository, userRepo repositories.UserRepository) WalletUseCaseInterface {
	return &WalletUseCase{
		walletRepo: walletRepo,
		userRepo:   userRepo,
	}
}

func (uc *WalletUseCase) CreateWallet(ctx context.Context, req *dto.CreateWalletRequest) (*dto.WalletResponse, error) {
	funcCtx := "CreateWallet"

	// Check if wallet already exists using name and user ID
	paramsGetOne := map[string]interface{}{"name": req.Name, "user_id": req.UserID}
	existingWallet, _ := uc.walletRepo.GetOne(ctx, paramsGetOne)
	if existingWallet != nil {
		logger.LogError(funcCtx, "wallet already exists", nil, logrus.Fields{"name": req.Name, "user_id": req.UserID})
		return nil, helpers.NewConflictError("wallet with this name already exists for this user", "")
	}

	// Create wallet entity
	wallet := &entities.Wallet{
		Name:     req.Name,
		Type:     req.Type,
		Category: req.Category,
		Balance:  req.Balance,
		Currency: req.Currency,
		UserID:   req.UserID,
	}

	// Save wallet
	if err := uc.walletRepo.Create(ctx, wallet); err != nil {
		logger.LogError(funcCtx, "failed to create wallet", err, logrus.Fields{"name": req.Name, "user_id": req.UserID})
		return nil, helpers.NewInternalError("failed to create wallet", err.Error())
	}

	return dto.MapToWalletResponse(wallet), nil
}

func (uc *WalletUseCase) GetWallet(ctx context.Context, id uuid.UUID, loggedUserID uuid.UUID) (*dto.WalletResponse, error) {
	funcCtx := "GetWallet"

	wallet, err := uc.walletRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get wallet", err, logrus.Fields{
			"wallet_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("wallet not found", "")
	}

	// Authorization check: non-admin users can only access their own wallets
	if loggedUserID != wallet.UserID {
		logger.LogError(funcCtx, "unauthorized access to wallet", nil, logrus.Fields{
			"wallet_id":      id.String(),
			"wallet_user_id": wallet.UserID.String(),
			"logged_user_id": loggedUserID.String(),
		})
		return nil, helpers.NewNotFoundError("wallet not found", "")
	}

	return dto.MapToWalletResponse(wallet), nil
}

func (uc *WalletUseCase) GetWallets(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.WalletResponse], error) {
	funcCtx := "GetWallets"

	wallets, err := uc.walletRepo.GetAll(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to get wallets", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to get wallets", err.Error())
	}

	total, err := uc.walletRepo.CountWithFilters(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to count wallets", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to count wallets", err.Error())
	}

	walletResponses := make([]dto.WalletResponse, len(wallets))
	for i, wallet := range wallets {
		walletResponses[i] = *dto.MapToWalletResponse(wallet)
	}

	paginationMeta := helpers.NewPaginationMeta(queryParams.Page, queryParams.Limit, total)

	return &dto.PaginationData[dto.WalletResponse]{
		Data: walletResponses,
		Meta: paginationMeta,
	}, nil
}

func (uc *WalletUseCase) UpdateWallet(ctx context.Context, id uuid.UUID, req *dto.UpdateWalletRequest) (*dto.WalletResponse, error) {
	funcCtx := "UpdateWallet"

	// Get existing wallet
	wallet, err := uc.walletRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get wallet", err, logrus.Fields{
			"wallet_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("wallet not found", "")
	}

	if req.Name != "" {
		wallet.Name = req.Name
	}
	if req.Type != "" {
		wallet.Type = req.Type
	}
	if req.Category != "" {
		wallet.Category = req.Category
	}
	if req.Balance >= 0 {
		wallet.Balance = req.Balance
	}
	if req.Currency != "" {
		wallet.Currency = req.Currency
	}

	if req.UserID != uuid.Nil {
		_, err := uc.userRepo.GetByID(ctx, req.UserID)
		if err != nil {
			logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
				"user_id": req.UserID.String(),
			})
			return nil, helpers.NewNotFoundError("user not found", "")
		}
		wallet.UserID = req.UserID
	}

	// Save updated wallet
	if err := uc.walletRepo.Update(ctx, wallet); err != nil {
		logger.LogError(funcCtx, "failed to update wallet", err, logrus.Fields{
			"wallet_id": id.String(),
		})
		return nil, helpers.NewInternalError("failed to update wallet", err.Error())
	}

	return dto.MapToWalletResponse(wallet), nil
}

func (uc *WalletUseCase) DeleteWallet(ctx context.Context, id uuid.UUID) error {
	funcCtx := "DeleteWallet"

	// Check if wallet exists
	_, err := uc.walletRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get wallet", err, logrus.Fields{
			"wallet_id": id.String(),
		})
		return helpers.NewNotFoundError("wallet not found", "")
	}

	// Use soft delete for default delete operation
	if err := uc.walletRepo.SoftDelete(ctx, id); err != nil {
		logger.LogError(funcCtx, "failed to delete wallet", err, logrus.Fields{
			"wallet_id": id.String(),
		})
		return helpers.NewInternalError("failed to delete wallet", err.Error())
	}

	return nil
}
