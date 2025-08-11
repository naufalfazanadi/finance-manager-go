package usecases

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/auth"
	"github.com/naufalfazanadi/finance-manager-go/pkg/encryption"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/minio"
	"github.com/naufalfazanadi/finance-manager-go/pkg/upload"
	"github.com/sirupsen/logrus"
)

type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest, profilePhoto *multipart.FileHeader) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest, profilePhoto *multipart.FileHeader) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error // This now does soft delete
	// Soft delete methods
	RestoreUser(ctx context.Context, id string) error
	GetUsersWithDeleted(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error)
	GetOnlyDeletedUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error)
	HardDeleteUser(ctx context.Context, id string) error // For permanent deletion
}

type UserUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserUseCaseInterface {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest, profilePhoto *multipart.FileHeader) (*dto.UserResponse, error) {
	funcCtx := "CreateUser"

	// Hash the email to check for existing user
	hashResult := encryption.HashSHA256(req.Email)
	if hashResult.Error != nil {
		logger.LogError(funcCtx, "failed to hash email", hashResult.Error, logrus.Fields{"email": req.Email})
		return nil, helpers.NewInternalError("failed to process email", hashResult.Error.Error())
	}

	emailHash := hashResult.Data.(string)

	// Check if user already exists using email hash
	existingUser, _ := uc.userRepo.GetByEmailHash(ctx, emailHash)
	if existingUser != nil {
		logger.LogError(funcCtx, "user already exists", nil, logrus.Fields{"email_hash": emailHash})
		return nil, helpers.NewConflictError("user with this email already exists", "")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash password", err, logrus.Fields{"email_hash": emailHash})
		return nil, helpers.NewInternalError("failed to hash password", err.Error())
	}

	// Set role to default user role (no longer accepting from request)
	role := entities.UserRoleUser

	// Upload profile photo first if provided
	var profilePhotoPath string
	if profilePhoto != nil {
		var err error
		profilePhotoPath, err = uc.uploadProfilePhoto(profilePhoto)
		if err != nil {
			logger.LogError(funcCtx, "failed to upload profile photo", err, logrus.Fields{"email": req.Email})
			return nil, helpers.NewInternalError("failed to upload profile photo", err.Error())
		}
	}

	// Create user entity
	user := &entities.User{
		Email:        req.Email,     // Set the plain email - it will be encrypted in BeforeCreate hook
		BirthDate:    req.BirthDate, // Set the birth date - it will be encrypted in BeforeCreate hook
		Name:         req.Name,
		Password:     hashedPassword,
		Role:         role,
		ProfilePhoto: profilePhotoPath,
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		// Revert profile photo upload if user creation fails
		if profilePhotoPath != "" {
			if revertErr := uc.revertProfilePhotoUpload(profilePhotoPath); revertErr != nil {
				logger.LogError(funcCtx, "failed to revert profile photo upload", revertErr, logrus.Fields{"photo_path": profilePhotoPath})
			}
		}
		logger.LogError(funcCtx, "failed to create user", err, logrus.Fields{"email_hash": emailHash})
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

func (uc *UserUseCase) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest, profilePhoto *multipart.FileHeader) (*dto.UserResponse, error) {
	funcCtx := "UpdateUser"

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

	// Upload profile photo if provided
	var oldProfilePhoto string
	var newProfilePhotoPath string
	if profilePhoto != nil {
		// Store old profile photo path for cleanup if needed
		oldProfilePhoto = user.ProfilePhoto

		var err error
		newProfilePhotoPath, err = uc.uploadProfilePhoto(profilePhoto)
		if err != nil {
			logger.LogError(funcCtx, "failed to upload profile photo", err, logrus.Fields{"user_id": id})
			return nil, helpers.NewInternalError("failed to upload profile photo", err.Error())
		}
		user.ProfilePhoto = newProfilePhotoPath
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}

	// Update birth date (can be set to nil to clear it)
	user.BirthDate = req.BirthDate

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		// Revert new profile photo upload if user update fails
		if newProfilePhotoPath != "" {
			if revertErr := uc.revertProfilePhotoUpload(newProfilePhotoPath); revertErr != nil {
				logger.LogError(funcCtx, "failed to revert profile photo upload", revertErr, logrus.Fields{"photo_path": newProfilePhotoPath})
			}
		}
		logger.LogError(funcCtx, "failed to update user", err, logrus.Fields{
			"user_id": id,
		})
		return nil, helpers.NewInternalError("failed to update user", err.Error())
	}

	// Clean up old profile photo after successful update (only if a new one was uploaded)
	if newProfilePhotoPath != "" && oldProfilePhoto != "" && oldProfilePhoto != newProfilePhotoPath {
		if cleanupErr := uc.revertProfilePhotoUpload(oldProfilePhoto); cleanupErr != nil {
			// Log the error but don't fail the request since the user update was successful
			logger.LogError(funcCtx, "failed to cleanup old profile photo", cleanupErr, logrus.Fields{
				"user_id":        id,
				"old_photo_path": oldProfilePhoto,
			})
		}
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

	// Use soft delete for default delete operation
	if err := uc.userRepo.SoftDelete(ctx, userID); err != nil {
		logger.LogError(funcCtx, "failed to delete user", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewInternalError("failed to delete user", err.Error())
	}

	return nil
}

func (uc *UserUseCase) mapToUserResponse(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Role:         string(user.Role),
		BirthDate:    user.BirthDate,
		Age:          user.GetAge(),
		ProfilePhoto: user.GetProfilePhotoURL(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

// RestoreUser restores a soft deleted user by ID
func (uc *UserUseCase) RestoreUser(ctx context.Context, id string) error {
	funcCtx := "RestoreUser"

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.LogError(funcCtx, "invalid user ID", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewValidationError("invalid user ID", "")
	}

	if err := uc.userRepo.Restore(ctx, userID); err != nil {
		logger.LogError(funcCtx, "failed to restore user", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewInternalError("failed to restore user", err.Error())
	}

	return nil
}

// GetUsersWithDeleted gets all users including soft deleted ones
func (uc *UserUseCase) GetUsersWithDeleted(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error) {
	funcCtx := "GetUsersWithDeleted"

	users, err := uc.userRepo.GetWithDeleted(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to get users with deleted", err, logrus.Fields{
			"limit":  queryParams.Limit,
			"offset": queryParams.GetOffset(),
		})
		return nil, helpers.NewInternalError("failed to get users", err.Error())
	}

	// Convert entities to DTOs
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *uc.mapToUserResponse(user))
	}

	// Get total count
	totalCount, err := uc.userRepo.CountWithFilters(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to count users", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to count users", err.Error())
	}

	totalPages := int((totalCount + int64(queryParams.Limit) - 1) / int64(queryParams.Limit))

	return &dto.UsersResponse{
		Users: userResponses,
		Pagination: &dto.PaginationMeta{
			Page:       queryParams.Page,
			Limit:      queryParams.Limit,
			Total:      totalCount,
			TotalPages: totalPages,
		},
	}, nil
}

// GetOnlyDeletedUsers gets only soft deleted users
func (uc *UserUseCase) GetOnlyDeletedUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.UsersResponse, error) {
	funcCtx := "GetOnlyDeletedUsers"

	users, err := uc.userRepo.GetOnlyDeleted(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to get deleted users", err, logrus.Fields{
			"limit":  queryParams.Limit,
			"offset": queryParams.GetOffset(),
		})
		return nil, helpers.NewInternalError("failed to get deleted users", err.Error())
	}

	// Convert entities to DTOs
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *uc.mapToUserResponse(user))
	}

	// For deleted users, we'll use a simple count since CountWithFilters doesn't handle deleted records
	totalCount := int64(len(users))
	totalPages := int((totalCount + int64(queryParams.Limit) - 1) / int64(queryParams.Limit))

	return &dto.UsersResponse{
		Users: userResponses,
		Pagination: &dto.PaginationMeta{
			Page:       queryParams.Page,
			Limit:      queryParams.Limit,
			Total:      totalCount,
			TotalPages: totalPages,
		},
	}, nil
}

// HardDeleteUser permanently deletes a user from the database
func (uc *UserUseCase) HardDeleteUser(ctx context.Context, id string) error {
	funcCtx := "HardDeleteUser"

	userID, err := uuid.Parse(id)
	if err != nil {
		logger.LogError(funcCtx, "invalid user ID", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewValidationError("invalid user ID", "")
	}

	if err := uc.userRepo.HardDelete(ctx, userID); err != nil {
		logger.LogError(funcCtx, "failed to hard delete user", err, logrus.Fields{
			"user_id": id,
		})
		return helpers.NewInternalError("failed to hard delete user", err.Error())
	}

	return nil
}

// uploadProfilePhoto handles the profile photo upload to minio
func (uc *UserUseCase) uploadProfilePhoto(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", nil
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Get file extension based on content type
	contentType := upload.GetFileExtension(upload.DetectContentType(fileContent))
	if contentType == "" {
		contentType = ".jpg" // fallback
	}

	// Generate unique filename with timestamp
	timestamp := time.Now().Unix()
	now := time.Now()
	folder := fmt.Sprintf("profile-photo/%d/%02d", now.Year(), now.Month())
	filename := fmt.Sprintf("profile_photo_%d%s", timestamp, contentType)

	// Initialize minio client
	minioClient, err := minio.NewClient()
	if err != nil {
		return "", fmt.Errorf("failed to initialize minio client: %w", err)
	}

	// Upload to minio (public bucket for profile photos)
	uploadResult, err := minioClient.UploadPublic(context.Background(), minio.UploadPublicDto{
		OriginalName: fileHeader.Filename,
		Folder:       folder,
		FileName:     filename,
		File:         fileContent,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to minio: %w", err)
	}

	return uploadResult.Path, nil
}

// Revert profile photo upload if user creation fails
func (uc *UserUseCase) revertProfilePhotoUpload(photoPath string) error {
	if photoPath == "" {
		return nil // Nothing to revert
	}

	minioClient, err := minio.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize minio client: %w", err)
	}

	if err := minioClient.RemoveObjectByPath(context.Background(), photoPath); err != nil {
		return fmt.Errorf("failed to delete profile photo from minio: %w", err)
	}

	return nil
}
