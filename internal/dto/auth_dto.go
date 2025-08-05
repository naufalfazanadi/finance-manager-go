package dto

import (
	"time"

	"github.com/google/uuid"
)

// Authentication Request DTOs
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
	Password string `json:"password" validate:"required,min=8,max=100" example:"password123"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin" example:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// Authentication Response DTOs
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// JWT Claims
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Exp    int64     `json:"exp"`
	Iat    int64     `json:"iat"`
}

// Token validation response
type TokenValidationResponse struct {
	Valid  bool      `json:"valid"`
	UserID uuid.UUID `json:"user_id,omitempty"`
	Email  string    `json:"email,omitempty"`
	Name   string    `json:"name,omitempty"`
	Role   string    `json:"role,omitempty"`
	Exp    time.Time `json:"exp,omitempty"`
}
