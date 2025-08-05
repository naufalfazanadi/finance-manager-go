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

type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserUseCase struct {
	userRepo  repositories.UserRepository
	validator *validator.Validator
}

func NewUserUseCase(userRepo repositories.UserRepository, validator *validator.Validator) UserUseCaseInterface {
	return &UserUseCase{
		userRepo:  userRepo,
		validator: validator,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	funcCtx := "CreateUser"

	// Validate request
	if err := uc.validator.Validate(req); err != nil {
		logger.LogError(funcCtx, "validation failed", err, logrus.Fields{"email": req.Email})
		return nil, helpers.NewValidationError("validation failed", err.Error())
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		logger.LogError(funcCtx, "user already exists", nil, logrus.Fields{"email": req.Email})
		return nil, helpers.NewConflictError("user with this email already exists", "")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash password", err, logrus.Fields{"email": req.Email})
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
		logger.LogError(funcCtx, "failed to create user", err, logrus.Fields{"email": req.Email})
		return nil, helpers.NewInternalError("failed to create user", err.Error())
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	funcCtx := "GetUser"

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.LogError(funcCtx, "invalid user ID format", err, logrus.Fields{"user_id": id})
		return nil, helpers.NewBadRequestError("invalid user ID format", "")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": id,
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error) {
	funcCtx := "GetUsers"

	users, err := uc.userRepo.GetAll(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to get users", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to get users", err.Error())
	}

	total, err := uc.userRepo.CountWithFilters(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to count users", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to count users", err.Error())
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *uc.mapToUserResponse(user)
	}

	paginationMeta := dto.NewPaginationMeta(queryParams.Page, queryParams.Limit, total)

	return &dto.UsersResponse{
		Users:      userResponses,
		Pagination: paginationMeta,
	}, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	funcCtx := "UpdateUser"

	// Validate request
	if err := uc.validator.Validate(req); err != nil {
		logger.LogError(funcCtx, "validation failed", err, logrus.Fields{
			"user_id": id,
		})
		return nil, helpers.NewValidationError("validation failed", err.Error())
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.LogError(funcCtx, "invalid user ID format", err, logrus.Fields{
			"user_id": id,
		})
		return nil, helpers.NewBadRequestError("invalid user ID format", "")
	}

	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": id,
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.LogError(funcCtx, "failed to update user", err, logrus.Fields{
			"user_id": id,
		})
		return nil, helpers.NewInternalError("failed to update user", err.Error())
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	funcCtx := "DeleteUser"

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.LogError(funcCtx, "invalid user ID format", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewBadRequestError("invalid user ID format", "")
	}

	// Check if user exists
	_, err = uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewNotFoundError("user not found", "")
	}

	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		logger.LogError(funcCtx, "failed to delete user", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewInternalError("failed to delete user", err.Error())
	}

	return nil
}

func (uc *UserUseCase) mapToUserResponse(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
