package helpers

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

// PaginatedSuccessResponse sends a successful response with pagination
func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, meta interface{}) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
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

// ErrorResponse sends an error response with consistent structure
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, details interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   details,
	})
}

// HandleError handles errors based on their type and sends appropriate response
func HandleError(c *fiber.Ctx, err error, defaultMessage string) error {
	if IsAppError(err) {
		appErr := err.(*AppError)
		return sendErrorByType(c, appErr.Type, appErr.Message, appErr.Details)
	}

	// Handle non-AppError types
	return InternalServerErrorResponse(c, defaultMessage, err.Error())
}

// sendErrorByType is a helper function that maps error types to HTTP status codes
func sendErrorByType(c *fiber.Ctx, errorType ErrorType, message string, details interface{}) error {
	statusCodeMap := map[ErrorType]int{
		ErrorTypeValidation:   fiber.StatusBadRequest,
		ErrorTypeNotFound:     fiber.StatusNotFound,
		ErrorTypeConflict:     fiber.StatusConflict,
		ErrorTypeBadRequest:   fiber.StatusBadRequest,
		ErrorTypeUnauthorized: fiber.StatusUnauthorized,
		ErrorTypeForbidden:    fiber.StatusForbidden,
		ErrorTypeInternal:     fiber.StatusInternalServerError,
	}

	statusCode, exists := statusCodeMap[errorType]
	if !exists {
		statusCode = fiber.StatusInternalServerError
	}

	return ErrorResponse(c, statusCode, message, details)
}

// Specific error response functions with uniform parameters
func BadRequestResponse(c *fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, details)
}

func UnauthorizedResponse(c *fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, details)
}

func ForbiddenResponse(c *fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, details)
}

func NotFoundResponse(c *fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, details)
}

func ConflictResponse(c *fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusConflict, message, details)
}

func InternalServerErrorResponse(c *fiber.Ctx, message string, details interface{}) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, details)
}

// NoContentResponse sends a no content response
func NoContentResponse(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
