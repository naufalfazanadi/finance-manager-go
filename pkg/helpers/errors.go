package helpers

import (
	"errors"
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict     ErrorType = "CONFLICT_ERROR"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST_ERROR"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED_ERROR"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN_ERROR"
)

// AppError represents a custom application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Error constructors
func NewValidationError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Details: details,
	}
}

func NewConflictError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Details: details,
	}
}

func NewBadRequestError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Details: details,
	}
}

func NewUnauthorizedError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Details: details,
	}
}

func NewInternalError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Details: details,
	}
}

func NewForbiddenError(message, details string) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Details: details,
	}
}

// GetErrorType returns the ErrorType from an error if it's an AppError
func GetErrorType(err error) ErrorType {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return ErrorTypeInternal
}

// IsAppError checks if the error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
