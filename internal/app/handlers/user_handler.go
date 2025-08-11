package handlers

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"

	"github.com/gofiber/fiber/v2"
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
// @Param profile_photo formData file false "Profile Photo (JPG/PNG, max 2MB)"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Handle form data with optional photo
	req := &dto.CreateUserRequest{}
	config := validator.FormValidationConfig{
		RequiredFields: []string{"email", "name", "password"},
		OptionalFields: []string{"birth_date"},
		FileFields:     []string{"profile_photo"},
		ValidateFiles:  true,
	}

	formResult, err := h.validator.ValidateForm(c, req, config)
	if err != nil {
		return helpers.HandleError(c, err, "Form validation failed")
	}

	// Get the validated DTO and profile photo
	createReq := formResult.Data.(*dto.CreateUserRequest)
	profilePhoto := formResult.Files["profile_photo"]

	// Create user with optional profile photo
	user, err := h.userUseCase.CreateUser(c.Context(), createReq, profilePhoto)
	if err != nil {
		return helpers.HandleError(c, err, "User creation failed")
	}

	return helpers.CreatedResponse(c, "User created successfully", user)
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
		return helpers.HandleError(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Validate ID format
	idParam := dto.IDParam{ID: id}
	if err := h.validator.Validate(&idParam); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid user ID format", err.Error()), "Invalid user ID format")
	}

	user, err := h.userUseCase.GetUser(c.Context(), id)
	if err != nil {
		return helpers.HandleError(c, err, "Failed to get user")
	}

	return helpers.SuccessResponse(c, "User retrieved successfully", user)
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
	queryParams := dto.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams.PaginationQuery); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid pagination parameters", err.Error()), "Invalid pagination parameters")
	}
	if err := h.validator.Validate(queryParams.FilterQuery); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid filter parameters", err.Error()), "Invalid filter parameters")
	}

	users, err := h.userUseCase.GetUsers(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleError(c, err, "Failed to get users")
	}

	return helpers.PaginatedSuccessResponse(c, "Users retrieved successfully", users.GetUsersData(), users.GetPaginationMeta())
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
// @Param profile_photo formData file false "Profile Photo (JPG/PNG, max 2MB)"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleError(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Validate ID format
	idParam := dto.IDParam{ID: id}
	if err := h.validator.Validate(&idParam); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid user ID format", err.Error()), "Invalid user ID format")
	}

	// Handle form data with optional photo
	req := &dto.UpdateUserRequest{}
	config := validator.FormValidationConfig{
		OptionalFields:   []string{"name", "birth_date"},
		FileFields:       []string{"profile_photo"},
		ValidateFiles:    true,
		AllowEmptyUpdate: true,
	}

	formResult, err := h.validator.ValidateForm(c, req, config)
	if err != nil {
		return helpers.HandleError(c, err, "Form validation failed")
	}

	// Get the validated DTO and profile photo
	updateReq := formResult.Data.(*dto.UpdateUserRequest)
	profilePhoto := formResult.Files["profile_photo"]

	// Update user with optional profile photo
	user, err := h.userUseCase.UpdateUser(c.Context(), id, updateReq, profilePhoto)
	if err != nil {
		return helpers.HandleError(c, err, "User update failed")
	}

	return helpers.SuccessResponse(c, "User updated successfully", user)
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
		return helpers.HandleError(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Validate ID format
	idParam := dto.IDParam{ID: id}
	if err := h.validator.Validate(&idParam); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid user ID format", err.Error()), "Invalid user ID format")
	}

	err := h.userUseCase.DeleteUser(c.Context(), id)
	if err != nil {
		return helpers.HandleError(c, err, "User deletion failed")
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
		return helpers.HandleError(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Validate ID format
	idParam := dto.IDParam{ID: id}
	if err := h.validator.Validate(&idParam); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid user ID format", err.Error()), "Invalid user ID format")
	}

	err := h.userUseCase.RestoreUser(c.Context(), id)
	if err != nil {
		return helpers.HandleError(c, err, "User restoration failed")
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
	queryParams := dto.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams.PaginationQuery); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid pagination parameters", err.Error()), "Invalid pagination parameters")
	}
	if err := h.validator.Validate(queryParams.FilterQuery); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid filter parameters", err.Error()), "Invalid filter parameters")
	}

	users, err := h.userUseCase.GetUsersWithDeleted(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleError(c, err, "Failed to get users with deleted")
	}

	return helpers.SuccessResponse(c, "Users retrieved successfully", users)
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
	queryParams := dto.ParseQueryParams(c)

	// Validate query parameters
	if err := h.validator.Validate(queryParams.PaginationQuery); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid pagination parameters", err.Error()), "Invalid pagination parameters")
	}
	if err := h.validator.Validate(queryParams.FilterQuery); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid filter parameters", err.Error()), "Invalid filter parameters")
	}

	users, err := h.userUseCase.GetOnlyDeletedUsers(c.Context(), queryParams)
	if err != nil {
		return helpers.HandleError(c, err, "Failed to get deleted users")
	}

	return helpers.SuccessResponse(c, "Deleted users retrieved successfully", users)
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
		return helpers.HandleError(c, helpers.NewBadRequestError("User ID is required", "ID parameter is missing"), "User ID is required")
	}

	// Validate ID format
	idParam := dto.IDParam{ID: id}
	if err := h.validator.Validate(&idParam); err != nil {
		return helpers.HandleError(c, helpers.NewValidationError("Invalid user ID format", err.Error()), "Invalid user ID format")
	}

	err := h.userUseCase.HardDeleteUser(c.Context(), id)
	if err != nil {
		return helpers.HandleError(c, err, "User hard deletion failed")
	}

	return helpers.NoContentResponse(c)
}
