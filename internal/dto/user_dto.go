package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateUserRequest struct {
	Email     string     `json:"email" validate:"required,email" example:"user@example.com"`
	Name      string     `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Password  string     `json:"password" validate:"required,min=8,max=100" example:"password123"`
	BirthDate *time.Time `json:"birth_date,omitempty" example:"1990-01-15T00:00:00Z"`
}

type UpdateUserRequest struct {
	Name      string     `json:"name" validate:"omitempty,min=2,max=100" example:"John Doe Updated"`
	BirthDate *time.Time `json:"birth_date,omitempty" example:"1990-01-15T00:00:00Z"`
}

// Response DTOs
type UserResponse struct {
	ID        uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string     `json:"email" example:"user@example.com"`
	Name      string     `json:"name" example:"John Doe"`
	Role      string     `json:"role" example:"user"`
	BirthDate *time.Time `json:"birth_date,omitempty" example:"1990-01-15T00:00:00Z"`
	Age       *int       `json:"age,omitempty" example:"33"`
	CreatedAt time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type UsersResponse struct {
	Users      []UserResponse  `json:"users"`
	Pagination *PaginationMeta `json:"pagination"`
}

// GetUsersData returns just the users array for the response
func (ur *UsersResponse) GetUsersData() []UserResponse {
	return ur.Users
}

// GetPaginationMeta returns just the pagination metadata
func (ur *UsersResponse) GetPaginationMeta() *PaginationMeta {
	return ur.Pagination
}

// Error Response
type ErrorResponse struct {
	Error   string                 `json:"error" example:"Bad Request"`
	Message string                 `json:"message" example:"Invalid input data"`
	Details map[string]interface{} `json:"details,omitempty"`
}
