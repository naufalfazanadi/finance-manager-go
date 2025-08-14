package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
)

// Request DTOs
type CreateWalletRequest struct {
	Name     string    `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Type     string    `json:"type" validate:"required" example:"personal"`
	Category string    `json:"category" validate:"required" example:"income"`
	Balance  float64   `json:"balance" validate:"omitempty,min=0" example:"1000.50"`
	Currency string    `json:"currency" validate:"omitempty" example:"IDR"`
	UserID   uuid.UUID `json:"user_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type UpdateWalletRequest struct {
	Name     string    `json:"name" validate:"omitempty,min=2,max=100" example:"John Doe Updated"`
	Type     string    `json:"type" validate:"omitempty" example:"personal"`
	Category string    `json:"category" validate:"omitempty" example:"income"`
	Balance  float64   `json:"balance" validate:"omitempty,min=0" example:"1000.50"`
	Currency string    `json:"currency" validate:"omitempty" example:"IDR"`
	UserID   uuid.UUID `json:"user_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Response DTOs
type WalletResponse struct {
	ID        uuid.UUID     `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string        `json:"name" example:"John Doe"`
	Type      string        `json:"type" example:"personal"`
	Category  string        `json:"category" example:"income"`
	Balance   float64       `json:"balance" example:"1000.50"`
	Currency  string        `json:"currency" example:"IDR"`
	UserID    uuid.UUID     `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt time.Time     `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time     `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// MapToWalletResponse converts a Wallet entity to WalletResponse DTO
func MapToWalletResponse(wallet *entities.Wallet) *WalletResponse {
	response := &WalletResponse{
		ID:        wallet.ID,
		Name:      wallet.Name,
		Type:      wallet.Type,
		Category:  wallet.Category,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		UserID:    wallet.UserID,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}

	// Include user data if it's preloaded
	if wallet.User.ID != uuid.Nil {
		response.User = MapToUserResponse(&wallet.User)
	}

	return response
}
