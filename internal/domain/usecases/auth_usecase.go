package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/auth"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"
	"github.com/sirupsen/logrus"
)

type AuthUseCaseInterface interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error)
}

type AuthUseCase struct {
	userRepo  repositories.UserRepository
	validator *validator.Validator
}

func NewAuthUseCase(userRepo repositories.UserRepository, validator *validator.Validator) AuthUseCaseInterface {
	return &AuthUseCase{
		userRepo:  userRepo,
		validator: validator,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	funcCtx := "Register"

	// Validate request
	if err := uc.validator.Validate(req); err != nil {
		logger.LogError(funcCtx, "validation failed", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewValidationError("validation failed", err.Error())
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		logger.LogError(funcCtx, "user already exists during registration", nil, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewConflictError("user with this email already exists", "")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash password", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewInternalError("failed to hash password", err.Error())
	}

	// Set role (default to user if not specified)
	role := entities.UserRoleUser
	if req.Role != "" {
		role = entities.UserRole(req.Role)
	}

	// Create user entity
	user := &entities.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
		Role:     role,
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.LogError(funcCtx, "failed to create user", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewInternalError("failed to create user", err.Error())
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.LogError(funcCtx, "failed to generate token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to generate token", err.Error())
	}

	logger.LogSuccess(funcCtx, "user registered successfully", logrus.Fields{
		"user_id": user.ID.String(),
		"email":   user.Email,
	})

	return &dto.AuthResponse{
		User:  *uc.mapToUserResponse(user),
		Token: token,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	funcCtx := "Login"

	// Validate request
	if err := uc.validator.Validate(req); err != nil {
		logger.LogError(funcCtx, "validation failed", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewValidationError("validation failed", err.Error())
	}

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user by email", err, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	// Check password
	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		logger.LogError(funcCtx, "login attempt with invalid password", nil, logrus.Fields{
			"email": req.Email,
		})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.LogError(funcCtx, "failed to generate token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to generate token", err.Error())
	}

	logger.LogSuccess(funcCtx, "user logged in successfully", logrus.Fields{
		"user_id": user.ID.String(),
		"email":   user.Email,
	})

	return &dto.LoginResponse{
		User:  *uc.mapToUserResponse(user),
		Token: token,
	}, nil
}

func (uc *AuthUseCase) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	funcCtx := "GetProfile"

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		logger.LogError(funcCtx, "invalid user ID format", err, logrus.Fields{
			"user_id": userID,
		})
		return nil, helpers.NewBadRequestError("invalid user ID format", "")
	}

	user, err := uc.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": userID,
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *AuthUseCase) mapToUserResponse(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
