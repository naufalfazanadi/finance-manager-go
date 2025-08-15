package usecases

import (
	"context"
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
	"github.com/sirupsen/logrus"
)

type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetUserWithPreload(ctx context.Context, id uuid.UUID, preloadRelations []string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.UserResponse], error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error // This now does soft delete
	// Soft delete methods
	RestoreUser(ctx context.Context, id uuid.UUID) error
	GetUsersWithDeleted(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.UserResponse], error)
	GetOnlyDeletedUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.UserResponse], error)
	HardDeleteUser(ctx context.Context, id uuid.UUID) error // For permanent deletion
}

type UserUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) UserUseCaseInterface {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
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
	if req.ProfilePhotoFile != nil {
		uploadResult, err := minio.UploadPhotoMinio(ctx, minio.UploadPhotoDto{
			FileHeader:   req.ProfilePhotoFile,
			FolderPrefix: "profile-photo",
			FilePrefix:   "profile_photo",
			BucketType:   minio.BucketTypePrivate, // Profile photos are private
		})
		if err != nil {
			logger.LogError(funcCtx, "failed to upload profile photo", err, logrus.Fields{"email": req.Email})
			return nil, helpers.NewInternalError("failed to upload profile photo", err.Error())
		}
		profilePhotoPath = uploadResult.Path
	}

	// Parse birth date if provided
	var birthDate *time.Time
	if req.BirthDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			logger.LogError(funcCtx, "failed to parse birth date", err, logrus.Fields{"birth_date": req.BirthDate})
			return nil, helpers.NewBadRequestError("invalid birth date format, expected YYYY-MM-DD", err.Error())
		}
		birthDate = &parsedTime
	}

	// Create user entity
	user := &entities.User{
		Email:        req.Email, // Set the plain email - it will be encrypted in BeforeCreate hook
		BirthDate:    birthDate, // Set the parsed birth date - it will be encrypted in BeforeCreate hook
		Name:         req.Name,
		Password:     hashedPassword,
		Role:         role,
		ProfilePhoto: profilePhotoPath,
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		// Revert profile photo upload if user creation fails
		if profilePhotoPath != "" {
			if revertErr := minio.DeletePhotoMinio(ctx, minio.DeletePhotoDto{
				PhotoPath:  profilePhotoPath,
				BucketType: minio.BucketTypePrivate, // Profile photos are in private bucket
			}); revertErr != nil {
				logger.LogError(funcCtx, "failed to revert profile photo upload", revertErr, logrus.Fields{"photo_path": profilePhotoPath})
			}
		}
		logger.LogError(funcCtx, "failed to create user", err, logrus.Fields{"email_hash": emailHash})
		return nil, helpers.NewInternalError("failed to create user", err.Error())
	}

	return dto.MapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	funcCtx := "GetUser"

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	return dto.MapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUserWithPreload(ctx context.Context, id uuid.UUID, preloadRelations []string) (*dto.UserResponse, error) {
	funcCtx := "GetUserWithPreload"

	user, err := uc.userRepo.GetByIDWithPreload(ctx, id, preloadRelations)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user with preload", err, logrus.Fields{
			"user_id": id.String(),
			"preload": preloadRelations,
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	return dto.MapToUserResponse(user), nil
}

func (uc *UserUseCase) GetUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.UserResponse], error) {
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
		userResponses[i] = *dto.MapToUserResponse(user)
	}

	paginationMeta := helpers.NewPaginationMeta(queryParams.Page, queryParams.Limit, total)

	return &dto.PaginationData[dto.UserResponse]{
		Data: userResponses,
		Meta: paginationMeta,
	}, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	funcCtx := "UpdateUser"

	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	// Upload profile photo if provided
	var oldProfilePhoto string
	var newProfilePhotoPath string
	if req.ProfilePhotoFile != nil {
		// Store old profile photo path for cleanup if needed
		oldProfilePhoto = user.ProfilePhoto

		uploadResult, err := minio.UploadPhotoMinio(ctx, minio.UploadPhotoDto{
			FileHeader:   req.ProfilePhotoFile,
			FolderPrefix: "profile-photo",
			FilePrefix:   "profile_photo",
			BucketType:   minio.BucketTypePrivate, // Profile photos are private
		})
		if err != nil {
			logger.LogError(funcCtx, "failed to upload profile photo", err, logrus.Fields{"user_id": id.String()})
			return nil, helpers.NewInternalError("failed to upload profile photo", err.Error())
		}
		newProfilePhotoPath = uploadResult.Path
		user.ProfilePhoto = newProfilePhotoPath
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}

	// Update birth date if provided
	if req.BirthDate != "" {
		parsedTime, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			// Revert new profile photo upload if date parsing fails
			if newProfilePhotoPath != "" {
				if revertErr := minio.DeletePhotoMinio(ctx, minio.DeletePhotoDto{
					PhotoPath:  newProfilePhotoPath,
					BucketType: minio.BucketTypePrivate,
				}); revertErr != nil {
					logger.LogError(funcCtx, "failed to revert profile photo upload", revertErr, logrus.Fields{"photo_path": newProfilePhotoPath})
				}
			}
			logger.LogError(funcCtx, "failed to parse birth date", err, logrus.Fields{"birth_date": req.BirthDate})
			return nil, helpers.NewBadRequestError("invalid birth date format, expected YYYY-MM-DD", err.Error())
		}
		user.BirthDate = &parsedTime
	}

	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		// Revert new profile photo upload if user update fails
		if newProfilePhotoPath != "" {
			if revertErr := minio.DeletePhotoMinio(ctx, minio.DeletePhotoDto{
				PhotoPath:  newProfilePhotoPath,
				BucketType: minio.BucketTypePrivate, // Profile photos are in public bucket
			}); revertErr != nil {
				logger.LogError(funcCtx, "failed to revert profile photo upload", revertErr, logrus.Fields{"photo_path": newProfilePhotoPath})
			}
		}
		logger.LogError(funcCtx, "failed to update user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return nil, helpers.NewInternalError("failed to update user", err.Error())
	}

	// Clean up old profile photo after successful update (only if a new one was uploaded)
	if newProfilePhotoPath != "" && oldProfilePhoto != "" && oldProfilePhoto != newProfilePhotoPath {
		if cleanupErr := minio.DeletePhotoMinio(ctx, minio.DeletePhotoDto{
			PhotoPath:  oldProfilePhoto,
			BucketType: minio.BucketTypePrivate, // Profile photos are in public bucket
		}); cleanupErr != nil {
			// Log the error but don't fail the request since the user update was successful
			logger.LogError(funcCtx, "failed to cleanup old profile photo", cleanupErr, logrus.Fields{
				"user_id":        id.String(),
				"old_photo_path": oldProfilePhoto,
			})
		}
	}

	return dto.MapToUserResponse(user), nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	funcCtx := "DeleteUser"

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return helpers.NewNotFoundError("user not found", "")
	}

	// Use soft delete for default delete operation
	if err := uc.userRepo.SoftDelete(ctx, id); err != nil {
		logger.LogError(funcCtx, "failed to delete user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return helpers.NewInternalError("failed to delete user", err.Error())
	}

	return nil
}

// RestoreUser restores a soft deleted user by ID
func (uc *UserUseCase) RestoreUser(ctx context.Context, id uuid.UUID) error {
	funcCtx := "RestoreUser"

	if err := uc.userRepo.Restore(ctx, id); err != nil {
		logger.LogError(funcCtx, "failed to restore user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return helpers.NewInternalError("failed to restore user", err.Error())
	}

	return nil
}

// GetUsersWithDeleted gets all users including soft deleted ones
func (uc *UserUseCase) GetUsersWithDeleted(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.UserResponse], error) {
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
		userResponses = append(userResponses, *dto.MapToUserResponse(user))
	}

	// Get total count
	totalCount, err := uc.userRepo.CountWithFilters(ctx, queryParams)
	if err != nil {
		logger.LogError(funcCtx, "failed to count users", err, logrus.Fields{})
		return nil, helpers.NewInternalError("failed to count users", err.Error())
	}

	totalPages := int((totalCount + int64(queryParams.Limit) - 1) / int64(queryParams.Limit))

	return &dto.PaginationData[dto.UserResponse]{
		Data: userResponses,
		Meta: &dto.PaginationMeta{
			Page:       queryParams.Page,
			Limit:      queryParams.Limit,
			Total:      totalCount,
			TotalPages: totalPages,
		},
	}, nil
}

// GetOnlyDeletedUsers gets only soft deleted users
func (uc *UserUseCase) GetOnlyDeletedUsers(ctx context.Context, queryParams *dto.QueryParams) (*dto.PaginationData[dto.UserResponse], error) {
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
		userResponses = append(userResponses, *dto.MapToUserResponse(user))
	}

	// For deleted users, we'll use a simple count since CountWithFilters doesn't handle deleted records
	totalCount := int64(len(users))
	totalPages := int((totalCount + int64(queryParams.Limit) - 1) / int64(queryParams.Limit))

	return &dto.PaginationData[dto.UserResponse]{
		Data: userResponses,
		Meta: &dto.PaginationMeta{
			Page:       queryParams.Page,
			Limit:      queryParams.Limit,
			Total:      totalCount,
			TotalPages: totalPages,
		},
	}, nil
}

// HardDeleteUser permanently deletes a user from the database
func (uc *UserUseCase) HardDeleteUser(ctx context.Context, id uuid.UUID) error {
	funcCtx := "HardDeleteUser"

	if err := uc.userRepo.HardDelete(ctx, id); err != nil {
		logger.LogError(funcCtx, "failed to hard delete user", err, logrus.Fields{
			"user_id": id.String(),
		})
		return helpers.NewInternalError("failed to hard delete user", err.Error())
	}

	return nil
}
