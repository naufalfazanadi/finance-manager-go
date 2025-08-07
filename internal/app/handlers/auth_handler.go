package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
)

type AuthHandler struct {
	authUseCase usecases.AuthUseCaseInterface
	validator   *validator.Validator
}

func NewAuthHandler(authUseCase usecases.AuthUseCaseInterface, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, name, password and birth date
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Use strict JSON parsing that rejects unknown fields
	if err := h.validator.ParseJSONStrict(c, &req); err != nil {
		return helpers.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequestResponse(c, "Validation failed", err.Error())
	}

	result, err := h.authUseCase.Register(c.Context(), &req)
	if err != nil {
		return helpers.HandleError(c, err, "Registration failed")
	}

	return helpers.CreatedResponse(c, "User registered successfully", result)
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := h.validator.ParseJSONStrict(c, &req); err != nil {
		return helpers.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequestResponse(c, "Validation failed", err.Error())
	}

	result, err := h.authUseCase.Login(c.Context(), &req)
	if err != nil {
		return helpers.HandleError(c, err, "Login failed")
	}

	return helpers.SuccessResponse(c, "Login successful", result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return helpers.UnauthorizedResponse(c, "Unauthorized", "User ID not found in token")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return helpers.UnauthorizedResponse(c, "Unauthorized", "Invalid user ID format in token")
	}

	// Validate user ID format
	idParam := dto.IDParam{ID: userIDStr}
	if err := h.validator.Validate(&idParam); err != nil {
		return helpers.UnauthorizedResponse(c, "Unauthorized", "Invalid user ID format in token")
	}

	result, err := h.authUseCase.GetProfile(c.Context(), userIDStr)
	if err != nil {
		return helpers.HandleError(c, err, "Failed to retrieve profile")
	}

	return helpers.SuccessResponse(c, "Profile retrieved successfully", result)
}
