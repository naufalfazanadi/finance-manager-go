package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Age   int    `json:"age" validate:"min=1,max=120"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
	Age  int    `json:"age" validate:"omitempty,min=1,max=120"`
}

// Response DTOs
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
