package handlers

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/upload"
	ut "github.com/naufalfazanadi/finance-manager-go/pkg/utils"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase usecases.UserUseCaseInterface
	validator   *validator.Validator
}

func NewUserHandler(userUseCase usecases.UserUseCaseInterface, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validator:   validator,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with email, name, password and optional profile photo using form data.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param name formData string true "Name"
// @Param password formData string true "Password"
// @Param birth_date formData string false "Birth Date (RFC3339 format)"
// @Param profile_photo_file formData file false "Profile Photo File (JPG/PNG, max 2MB)"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	// Parse and validate form data with strict field validation
	if err := h.validator.ParseFormAndValidate(c, &req); err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(fiberErr.Message, fiberErr.Error()), fiberErr.Message)
		}
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Validation failed", err.Error()), "Validation failed")
	}

	// Validate profile photo file if uploaded
	profilePhotoResult := h.validator.ValidateFile(req.ProfilePhotoFile, upload.ProfilePhotoValidation)
	if !profilePhotoResult.Valid {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("Profile photo validation failed", profilePhotoResult.Error), "Profile photo validation failed")
	}

	// Create user
	user, err := h.userUseCase.CreateUser(c.Context(), &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedCreateMsg("User"))
	}

	return helpers.CreatedResponse(c, ut.SuccessCreateMsg("User"), user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	userID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	if userID != c.Locals("userID").(uuid.UUID) && c.Locals("userRole") != "admin" {
		return helpers.HandleErrorResponse(c, helpers.NewForbiddenError("You do not have permission", "Permission denied"), "Permission denied")
	}

	user, err := h.userUseCase.GetUser(c.Context(), userID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("User"))
	}

	return helpers.SuccessResponse(c, ut.SuccessRetrieveMsg("User"), user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with pagination and filtering. Returns users array in data field and pagination info in meta field.
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term for name or email"
// @Param sort_by query string false "Field to sort by (name, email, created_at)"
// @Param sort_dir query string false "Sort direction (asc, desc)" default(asc)
// @Param name query string false "Filter by exact name"
// @Success 200 {object} object{success=bool,message=string,data=[]dto.UserResponse,meta=dto.PaginationMeta}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	queryParams := helpers.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError(ut.MsgErrInvalidQueryParams, err.Error()), ut.MsgErrInvalidQueryParams)
	}

	users, err := h.userUseCase.GetUsers(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Users"))
	}

	return helpers.PaginatedSuccessResponse(c, ut.SuccessRetrieveMsg("Users"), users.Data, users.Meta)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information with optional profile photo using form data.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "User ID"
// @Param name formData string false "Name"
// @Param birth_date formData string false "Birth Date (RFC3339 format)"
// @Param profile_photo_file formData file false "Profile Photo File (JPG/PNG, max 2MB)"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	userID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	var req dto.UpdateUserRequest

	// Parse and validate form data with strict field validation
	if err := h.validator.ParseFormAndValidate(c, &req); err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(fiberErr.Message, fiberErr.Error()), fiberErr.Message)
		}
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Validation failed", err.Error()), "Validation failed")
	}

	// Validate profile photo file if uploaded
	profilePhotoResult := h.validator.ValidateFile(req.ProfilePhotoFile, upload.ProfilePhotoValidation)
	if !profilePhotoResult.Valid {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("Profile photo validation failed", profilePhotoResult.Error), "Profile photo validation failed")
	}

	// Update user
	user, err := h.userUseCase.UpdateUser(c.Context(), userID, &req)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedUpdateMsg("User"))
	}

	return helpers.SuccessResponse(c, ut.SuccessUpdateMsg("User"), user)
}

// DeleteUser godoc
// @Summary Delete user by ID (soft delete)
// @Description Soft delete a user by their ID (marks as deleted but keeps data)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Parse and validate UUID format
	userID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("Invalid user ID format", "User ID must be a valid UUID"), "Invalid user ID format")
	}

	err = h.userUseCase.DeleteUser(c.Context(), userID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedDeleteMsg("User"))
	}

	return helpers.NoContentResponse(c)
}

// RestoreUser godoc
// @Summary Restore soft deleted user by ID
// @Description Restore a soft deleted user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{id}/restore [patch]
func (h *UserHandler) RestoreUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Parse and validate UUID format
	userID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("Invalid user ID format", "User ID must be a valid UUID"), "Invalid user ID format")
	}

	err = h.userUseCase.RestoreUser(c.Context(), userID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedRestoreMsg("User"))
	}

	return helpers.NoContentResponse(c)
}

// GetUsersWithDeleted godoc
// @Summary Get all users including soft deleted ones
// @Description Get all users including those that have been soft deleted
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field"
// @Param sort_type query string false "Sort type (asc/desc)" default(desc)
// @Success 200 {object} dto.UsersResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/with-deleted [get]
func (h *UserHandler) GetUsersWithDeleted(c *fiber.Ctx) error {
	queryParams := helpers.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams.PaginationQuery); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Invalid pagination parameters", err.Error()), "Invalid pagination parameters")
	}
	if err := h.validator.Validate(queryParams.FilterQuery); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Invalid filter parameters", err.Error()), "Invalid filter parameters")
	}

	users, err := h.userUseCase.GetUsersWithDeleted(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Users with deleted"))
	}

	return helpers.SuccessResponse(c, ut.SuccessRetrieveMsg("Users with deleted"), users)
}

// GetOnlyDeletedUsers godoc
// @Summary Get only soft deleted users
// @Description Get only users that have been soft deleted
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort by field"
// @Param sort_type query string false "Sort type (asc/desc)" default(desc)
// @Success 200 {object} dto.UsersResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/deleted [get]
func (h *UserHandler) GetOnlyDeletedUsers(c *fiber.Ctx) error {
	queryParams := helpers.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams.PaginationQuery); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Invalid pagination parameters", err.Error()), "Invalid pagination parameters")
	}
	if err := h.validator.Validate(queryParams.FilterQuery); err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewValidationError("Invalid filter parameters", err.Error()), "Invalid filter parameters")
	}

	users, err := h.userUseCase.GetOnlyDeletedUsers(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Deleted users"))
	}

	return helpers.SuccessResponse(c, ut.SuccessRetrieveMsg("Deleted users"), users)
}

// HardDeleteUser godoc
// @Summary Permanently delete user by ID
// @Description Permanently delete a user by their ID (cannot be restored)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{id}/hard-delete [delete]
func (h *UserHandler) HardDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Parse and validate UUID format
	userID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError("Invalid user ID format", "User ID must be a valid UUID"), "Invalid user ID format")
	}

	err = h.userUseCase.HardDeleteUser(c.Context(), userID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedDeleteMsg("User"))
	}

	return helpers.NoContentResponse(c)
}
