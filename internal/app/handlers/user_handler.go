package handlers

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUseCase usecases.UserUseCaseInterface
}

func NewUserHandler(userUseCase usecases.UserUseCaseInterface) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with email, name, and age
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	user, err := h.userUseCase.CreateUser(c.Context(), &req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return utils.ConflictResponse(c, "User creation failed", err.Error())
		} else if len(err.Error()) >= 10 && err.Error()[:10] == "validation" {
			return utils.BadRequestResponse(c, "Validation failed", err.Error())
		}
		return utils.InternalServerErrorResponse(c, "User creation failed", err.Error())
	}

	return utils.CreatedResponse(c, "User created successfully", user)
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
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequestResponse(c, "User ID is required", nil)
	}

	user, err := h.userUseCase.GetUser(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			return utils.NotFoundResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "Failed to get user", err.Error())
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
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
// @Param sort_by query string false "Field to sort by (name, email, age, created_at)"
// @Param sort_dir query string false "Sort direction (asc, desc)" default(asc)
// @Param name query string false "Filter by exact name"
// @Param age query int false "Filter by exact age"
// @Success 200 {object} object{success=bool,message=string,data=[]dto.UserResponse,meta=dto.PaginationMeta}
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	queryParams := dto.ParseQueryParams(c)

	users, err := h.userUseCase.GetUsers(c.Context(), queryParams)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to get users", err.Error())
	}

	return utils.PaginatedSuccessResponse(c, "Users retrieved successfully", users.GetUsersData(), users.GetPaginationMeta())
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequestResponse(c, "User ID is required", nil)
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	user, err := h.userUseCase.UpdateUser(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			return utils.NotFoundResponse(c, err.Error())
		} else if len(err.Error()) >= 10 && err.Error()[:10] == "validation" {
			return utils.BadRequestResponse(c, "Validation failed", err.Error())
		}
		return utils.InternalServerErrorResponse(c, "User update failed", err.Error())
	}

	return utils.SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequestResponse(c, "User ID is required", nil)
	}

	err := h.userUseCase.DeleteUser(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
			return utils.NotFoundResponse(c, err.Error())
		}
		return utils.InternalServerErrorResponse(c, "User deletion failed", err.Error())
	}

	return utils.NoContentResponse(c)
}
