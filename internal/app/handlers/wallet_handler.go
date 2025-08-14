package handlers

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	ut "github.com/naufalfazanadi/finance-manager-go/pkg/utils"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletUseCase usecases.WalletUseCaseInterface
	validator     *validator.Validator
}

func NewWalletHandler(walletUseCase usecases.WalletUseCaseInterface, validator *validator.Validator) *WalletHandler {
	return &WalletHandler{
		walletUseCase: walletUseCase,
		validator:     validator,
	}
}

func (h *WalletHandler) CreateWallet(c *fiber.Ctx) error {
	var req dto.CreateWalletRequest

	// Apply business logic for user authorization
	if c.Locals("userRole") != "admin" && req.UserID == uuid.Nil {
		req.UserID = c.Locals("userID").(uuid.UUID)
	}
	if req.UserID != c.Locals("userID").(uuid.UUID) && c.Locals("userRole") != "admin" {
		return helpers.HandleErrorResponse(c, helpers.NewForbiddenError("You do not have permission to create a wallet for this user", "Permission denied"), "Permission denied")
	}

	// Parse strict JSON validation and struct validation
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrReqBody, err.Error()), ut.MsgErrReqBody)
	}

	wallet, err := h.walletUseCase.CreateWallet(c.Context(), &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedCreateMsg("Wallet"))
	}

	return helpers.CreatedResponse(c, ut.SuccessCreateMsg("Wallet"), wallet)
}

func (h *WalletHandler) GetWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.ErrIDRequired, ut.MsgErrIDRequired), ut.ErrIDRequired)
	}

	// Parse and validate UUID format
	walletID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	// Get logged user information from context
	var loggedUserID uuid.UUID
	if c.Locals("userRole") != "admin" {
		loggedUserID = c.Locals("userID").(uuid.UUID)
	}

	wallet, err := h.walletUseCase.GetWallet(c.Context(), walletID, loggedUserID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Wallet"))
	}

	return helpers.SuccessResponse(c, ut.SuccessRetrieveMsg("Wallet"), wallet)
}

func (h *WalletHandler) GetWallets(c *fiber.Ctx) error {
	queryParams := helpers.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrInvalidQueryParams, err.Error()), ut.MsgErrInvalidQueryParams)
	}

	if c.Locals("userRole") != "admin" {
		queryParams.LoggedUserID = c.Locals("userID").(uuid.UUID)
	}

	wallets, err := h.walletUseCase.GetWallets(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Wallets"))
	}

	return helpers.PaginatedSuccessResponse(c, ut.SuccessRetrieveMsg("Wallets"), wallets.Data, wallets.Meta)
}

func (h *WalletHandler) UpdateWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	walletID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	var req dto.UpdateWalletRequest

	// Parse strict JSON validation and struct validation
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrReqBody, err.Error()), ut.MsgErrReqBody)
	}

	wallet, err := h.walletUseCase.UpdateWallet(c.Context(), walletID, &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedUpdateMsg("Wallet"))
	}

	return helpers.SuccessResponse(c, ut.SuccessUpdateMsg("Wallet"), wallet)
}

func (h *WalletHandler) DeleteWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	walletID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	err = h.walletUseCase.DeleteWallet(c.Context(), walletID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedDeleteMsg("Wallet"))
	}

	return helpers.NoContentResponse(c)
}
