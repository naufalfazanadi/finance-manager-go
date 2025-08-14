package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/auth"
	"github.com/naufalfazanadi/finance-manager-go/pkg/encryption"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

type AuthUseCaseInterface interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}

type AuthUseCase struct {
	userRepo repositories.UserRepository
}

func NewAuthUseCase(userRepo repositories.UserRepository) AuthUseCaseInterface {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	funcCtx := "Register"

	// Validate password strength
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		logger.LogError(funcCtx, "password validation failed", err, logrus.Fields{"email": req.Email})
		return nil, helpers.NewBadRequestError("password validation failed", err.Error())
	}

	// Hash the email to check for existing user
	hashResult := encryption.HashSHA256(req.Email)
	if hashResult.Error != nil {
		logger.LogError(funcCtx, "failed to hash email", hashResult.Error, logrus.Fields{"email": req.Email})
		return nil, helpers.NewInternalError("failed to process email", hashResult.Error.Error())
	}

	emailHash := hashResult.Data.(string)

	// Check if user already exists using email hash
	existingUser, errUser := uc.userRepo.GetByEmailHash(ctx, emailHash)
	if existingUser != nil {
		logger.LogError(funcCtx, "user already exists during registration", nil, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewConflictError("user with this email already exists", "")
	}
	fmt.Println(errUser)

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash password", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewInternalError("failed to hash password", err.Error())
	}

	// Parse birth date from string to time.Time
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		logger.LogError(funcCtx, "failed to parse birth date", err, logrus.Fields{
			"birth_date": req.BirthDate,
		})
		return nil, helpers.NewBadRequestError("invalid birth date format", "birth date must be in YYYY-MM-DD format")
	}

	// Set role to default user role (no longer accepting from request)
	role := entities.UserRoleUser

	// Create user entity
	user := &entities.User{
		Email:     req.Email,  // Set the plain email - it will be encrypted in BeforeCreate hook
		BirthDate: &birthDate, // Set the parsed birth date - it will be encrypted in BeforeCreate hook
		Name:      req.Name,
		Password:  hashedPassword,
		Role:      role,
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.LogError(funcCtx, "failed to create user", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewInternalError("failed to create user", err.Error())
	}

	// Generate JWT token (user.Email should be decrypted by AfterFind hook)
	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.LogError(funcCtx, "failed to generate token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to generate token", err.Error())
	}

	return &dto.AuthResponse{
		UserResponse: *dto.MapToUserResponse(user),
		Token:        token,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	funcCtx := "Login"

	// Hash the email to lookup user
	hashResult := encryption.HashSHA256(req.Email)
	if hashResult.Error != nil {
		logger.LogError(funcCtx, "failed to hash email", hashResult.Error, logrus.Fields{"email": req.Email})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	emailHash := hashResult.Data.(string)

	// Get user by email hash
	user, err := uc.userRepo.GetByEmailHash(ctx, emailHash)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user by email hash", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	// Check password
	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		logger.LogError(funcCtx, "login attempt with invalid password", nil, logrus.Fields{
			"email_hash": emailHash,
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

	return &dto.LoginResponse{
		UserResponse: *dto.MapToUserResponse(user),
		Token:        token,
	}, nil
}

func (uc *AuthUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	funcCtx := "GetProfile"

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	return dto.MapToUserResponse(user), nil
}
