package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PaginatedResponse represents the standard API response structure with pagination
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse sends a created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequestResponse sends a bad request error response
func BadRequestResponse(c *fiber.Ctx, message string, err interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, err)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, nil)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, nil)
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *fiber.Ctx, message string, err interface{}) error {
	return ErrorResponse(c, fiber.StatusConflict, message, err)
}

// UnprocessableEntityResponse sends an unprocessable entity error response
func UnprocessableEntityResponse(c *fiber.Ctx, message string, err interface{}) error {
	return ErrorResponse(c, fiber.StatusUnprocessableEntity, message, err)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *fiber.Ctx, message string, err interface{}) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, err)
}

// PaginatedSuccessResponse sends a successful response with pagination
func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, meta interface{}) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Success: false,
		Message: "Validation failed",
		Error:   errors,
	})
}

// NoContentResponse sends a no content response
func NoContentResponse(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
