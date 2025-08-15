package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
)

// Request DTOs
type CreateTransactionRequest struct {
	Name      string    `json:"name" validate:"required,min=2,max=255" example:"Grocery Shopping"`
	Cost      float64   `json:"cost" validate:"required,min=0" example:"50000.00"`
	Type      string    `json:"type" validate:"required,oneof=income expense" example:"expense"`
	Note      string    `json:"note" validate:"omitempty,max=1000" example:"Weekly grocery shopping at supermarket"`
	TCategory string    `json:"t_category" validate:"required,min=2,max=100" example:"food"`
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	WalletID  uuid.UUID `json:"wallet_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
}

type UpdateTransactionRequest struct {
	Name      string    `json:"name" validate:"omitempty,min=2,max=255" example:"Updated Grocery Shopping"`
	Cost      float64   `json:"cost" validate:"omitempty,min=0" example:"45000.00"`
	Type      string    `json:"type" validate:"omitempty,oneof=income expense" example:"expense"`
	Note      string    `json:"note" validate:"omitempty,max=1000" example:"Updated note for grocery shopping"`
	TCategory string    `json:"t_category" validate:"omitempty,min=2,max=100" example:"food"`
	UserID    uuid.UUID `json:"user_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	WalletID  uuid.UUID `json:"wallet_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
}

// Response DTOs
type TransactionResponse struct {
	ID        uuid.UUID       `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string          `json:"name" example:"Grocery Shopping"`
	Cost      float64         `json:"cost" example:"50000.00"`
	Type      string          `json:"type" example:"expense"`
	Note      string          `json:"note" example:"Weekly grocery shopping at supermarket"`
	TCategory string          `json:"t_category" example:"food"`
	UserID    uuid.UUID       `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	WalletID  uuid.UUID       `json:"wallet_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	User      *UserResponse   `json:"user,omitempty"`
	Wallet    *WalletResponse `json:"wallet,omitempty"`
	CreatedAt time.Time       `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time       `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// MapToTransactionResponse converts a Transaction entity to TransactionResponse DTO
func MapToTransactionResponse(transaction *entities.Transaction) *TransactionResponse {
	response := &TransactionResponse{
		ID:        transaction.ID,
		Name:      transaction.Name,
		Cost:      transaction.Cost,
		Type:      string(transaction.Type),
		Note:      transaction.Note,
		TCategory: transaction.TCategory,
		UserID:    transaction.UserID,
		WalletID:  transaction.WalletID,
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}

	// Include user data if it's preloaded
	if transaction.User.ID != uuid.Nil {
		response.User = MapToUserResponse(&transaction.User)
	}

	// Include wallet data if it's preloaded
	if transaction.Wallet.ID != uuid.Nil {
		response.Wallet = MapToWalletResponse(&transaction.Wallet)
	}

	return response
}
