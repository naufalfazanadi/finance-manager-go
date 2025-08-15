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

type TransactionHandler struct {
	transactionUseCase usecases.TransactionUseCaseInterface
	validator          *validator.Validator
}

func NewTransactionHandler(transactionUseCase usecases.TransactionUseCaseInterface, validator *validator.Validator) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
		validator:          validator,
	}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req dto.CreateTransactionRequest

	// Apply business logic for user authorization
	if c.Locals("userRole") != "admin" && req.UserID == uuid.Nil {
		req.UserID = c.Locals("userID").(uuid.UUID)
	}
	if req.UserID != c.Locals("userID").(uuid.UUID) && c.Locals("userRole") != "admin" {
		return helpers.HandleErrorResponse(c, helpers.NewForbiddenError("You do not have permission to create a transaction for this user", "Permission denied"), "Permission denied")
	}

	// Parse strict JSON validation and struct validation
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrReqBody, err.Error()), ut.MsgErrReqBody)
	}

	transaction, err := h.transactionUseCase.CreateTransaction(c.Context(), &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedCreateMsg("Transaction"))
	}

	return helpers.CreatedResponse(c, ut.SuccessCreateMsg("Transaction"), transaction)
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.ErrIDRequired, ut.MsgErrIDRequired), ut.ErrIDRequired)
	}

	// Parse and validate UUID format
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	// Get logged user information from context
	var loggedUserID uuid.UUID
	if c.Locals("userRole") != "admin" {
		loggedUserID = c.Locals("userID").(uuid.UUID)
	}

	transaction, err := h.transactionUseCase.GetTransaction(c.Context(), transactionID, loggedUserID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Transaction"))
	}

	return helpers.SuccessResponse(c, ut.SuccessRetrieveMsg("Transaction"), transaction)
}

func (h *TransactionHandler) GetTransactions(c *fiber.Ctx) error {
	queryParams := helpers.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrInvalidQueryParams, err.Error()), ut.MsgErrInvalidQueryParams)
	}

	if c.Locals("userRole") != "admin" {
		queryParams.LoggedUserID = c.Locals("userID").(uuid.UUID)
	}

	transactions, err := h.transactionUseCase.GetTransactions(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Transactions"))
	}

	return helpers.PaginatedSuccessResponse(c, ut.SuccessRetrieveMsg("Transactions"), transactions.Data, transactions.Meta)
}

func (h *TransactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	var req dto.UpdateTransactionRequest

	// Parse strict JSON validation and struct validation
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrReqBody, err.Error()), ut.MsgErrReqBody)
	}

	transaction, err := h.transactionUseCase.UpdateTransaction(c.Context(), transactionID, &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedUpdateMsg("Transaction"))
	}

	return helpers.SuccessResponse(c, ut.SuccessUpdateMsg("Transaction"), transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	err = h.transactionUseCase.DeleteTransaction(c.Context(), transactionID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedDeleteMsg("Transaction"))
	}

	return helpers.NoContentResponse(c)
}
