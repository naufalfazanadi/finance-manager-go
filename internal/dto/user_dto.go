package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
)

// Request DTOs
type CreateUserRequest struct {
	Email            string                `json:"email" form:"email" validate:"required,email" example:"user@example.com"`
	Name             string                `json:"name" form:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Password         string                `json:"password" form:"password" validate:"required,min=8,max=100" example:"password123"`
	BirthDate        string                `json:"birth_date,omitempty" form:"birth_date" validate:"omitempty,datetime=2006-01-02" example:"1990-01-15"`
	ProfilePhoto     string                `json:"profile_photo,omitempty" form:"profile_photo" validate:"omitempty" example:"profile photo URL or path"`
	ProfilePhotoFile *multipart.FileHeader `json:"-" form:"profile_photo_file" validate:"omitempty" swaggerignore:"true"`
}

type UpdateUserRequest struct {
	Name             string                `json:"name" form:"name" validate:"omitempty,min=2,max=100" example:"John Doe Updated"`
	BirthDate        string                `json:"birth_date,omitempty" form:"birth_date" validate:"omitempty,datetime=2006-01-02" example:"1990-01-15"`
	ProfilePhoto     string                `json:"profile_photo,omitempty" form:"profile_photo" validate:"omitempty" example:"profile photo URL or path"`
	ProfilePhotoFile *multipart.FileHeader `json:"-" form:"profile_photo_file" validate:"omitempty" swaggerignore:"true"`
}

// Response DTOs
type UserResponse struct {
	ID           uuid.UUID             `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email        string                `json:"email" example:"user@example.com"`
	Name         string                `json:"name" example:"John Doe"`
	Role         string                `json:"role" example:"user"`
	BirthDate    *time.Time            `json:"birth_date" example:"1990-01-15"`
	Age          *int                  `json:"age" example:"33"`
	ProfilePhoto string                `json:"profile_photo" example:"https://minio.example.com/public/profile-photo/2023/01/profile_photo_1641024000.jpg"`
	CreatedAt    time.Time             `json:"created_at" example:"2023-01-01"`
	UpdatedAt    time.Time             `json:"updated_at" example:"2023-01-01"`
	Wallets      []WalletResponse      `json:"wallets,omitempty"`
	Transactions []TransactionResponse `json:"transactions,omitempty"`
}

// MapToUserResponse converts a User entity to UserResponse DTO
func MapToUserResponse(user *entities.User) *UserResponse {
	response := &UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         string(user.Role),
		BirthDate:    user.BirthDate,
		Age:          user.GetAge(),
		ProfilePhoto: user.GetProfilePhotoURL(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	// Include wallets data if it's preloaded
	if len(user.Wallets) > 0 {
		response.Wallets = make([]WalletResponse, len(user.Wallets))
		for i, wallet := range user.Wallets {
			response.Wallets[i] = *MapToWalletResponse(&wallet)
		}
	}

	// Include transactions data if it's preloaded
	if len(user.Transactions) > 0 {
		response.Transactions = make([]TransactionResponse, len(user.Transactions))
		for i, transaction := range user.Transactions {
			response.Transactions[i] = *MapToTransactionResponse(&transaction)
		}
	}

	return response
}
