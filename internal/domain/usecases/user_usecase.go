package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"

	"github.com/google/uuid"
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
	// Validate request
	if err := uc.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create user entity
	user := &entities.User{
		Email: req.Email,
		Name:  req.Name,
		Age:   req.Age,
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error) {
	// TODO: Implement filtering in repository layer
	// For now, we're using basic pagination. To implement filtering:
	// 1. Add GetAllWithFilters method to UserRepository interface
	// 2. Pass queryParams.FilterQuery to the repository method
	// 3. Implement filtering logic in the infrastructure layer

	users, err := uc.userRepo.GetAll(ctx, queryParams.Limit, queryParams.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
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
	// Validate request
	if err := uc.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return uc.mapToUserResponse(user), nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Check if user exists
	_, err = uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (uc *UserUseCase) mapToUserResponse(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
