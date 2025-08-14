package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	// Parse and validate request
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Validation failed", err.Error()), "Invalid request body")
	}

	result, err := h.authUseCase.Register(c.Context(), &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, "Registration failed")
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

	// Parse and validate request
	if err := h.validator.ParseAndValidate(c, &req); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Validation failed", err.Error()), "Invalid request body")
	}

	result, err := h.authUseCase.Login(c.Context(), &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, "Login failed")
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
	userID := c.Locals("userID").(uuid.UUID)
	if userID == uuid.Nil {
		return helpers.HandleErrorResponse(c, helpers.NewUnauthorizedError("Unauthorized", "User ID not found in token"), "Unauthorized")
	}

	result, err := h.authUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, "Failed to retrieve profile")
	}

	return helpers.SuccessResponse(c, "Profile retrieved successfully", result)
}
