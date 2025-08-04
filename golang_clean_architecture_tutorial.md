# Go Clean Architecture Tutorial with Fiber, GORM & PostgreSQL

This tutorial will guide you through creating a Go backend with Clean Architecture, featuring User CRUD operations using Fiber (web framework), GORM (ORM), and PostgreSQL.

## Prerequisites

- Go 1.21+ installed
- PostgreSQL installed and running
- Basic understanding of Go programming
- Docker (optional, for PostgreSQL)

## Step 1: Project Setup

### 1.1 Create Project Directory
```bash
mkdir github.com/naufalfazanadi/finance-manager-go
cd github.com/naufalfazanadi/finance-manager-go
```

### 1.2 Initialize Go Module
```bash
go mod init github.com/naufalfazanadi/finance-manager-go
```

### 1.3 Create Project Structure
```bash
# Create directories
mkdir -p cmd/server
mkdir -p internal/{config,domain/{entities,repositories,services},infrastructure/{database,repositories},application/{usecases,dto},interfaces/http/{handlers,routes,middleware}}
mkdir -p pkg/{logger,validator,utils}
mkdir -p scripts
mkdir -p docs

# Create main files
touch cmd/server/main.go
touch internal/config/config.go
touch .env
touch README.md
```

## Step 2: Install Dependencies

```bash
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10
go get github.com/google/uuid
go get go.uber.org/zap
```

## Step 3: Configuration Setup

### 3.1 Environment Variables (.env)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=clean_api_db
DB_SSLMODE=disable

# Server
SERVER_PORT=8080
SERVER_HOST=localhost

# Application
APP_ENV=development
LOG_LEVEL=debug
```

### 3.2 Configuration Structure (internal/config/config.go)
```go
package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    App      AppConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type ServerConfig struct {
    Host string
    Port string
}

type AppConfig struct {
    Env      string
    LogLevel string
}

func LoadConfig() *Config {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "password"),
            DBName:   getEnv("DB_NAME", "clean_api_db"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Server: ServerConfig{
            Host: getEnv("SERVER_HOST", "localhost"),
            Port: getEnv("SERVER_PORT", "8080"),
        },
        App: AppConfig{
            Env:      getEnv("APP_ENV", "development"),
            LogLevel: getEnv("LOG_LEVEL", "debug"),
        },
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

## Step 4: Domain Layer

### 4.1 User Entity (internal/domain/entities/user.go)
```go
package entities

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null"`
    Name      string         `json:"name" gorm:"not null"`
    Age       int            `json:"age"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

// TableName sets the table name
func (User) TableName() string {
    return "users"
}
```

### 4.2 Repository Interface (internal/domain/repositories/user_repository.go)
```go
package repositories

import (
    "context"
    "github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
    "github.com/google/uuid"
)

type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    GetByEmail(ctx context.Context, email string) (*entities.User, error)
    GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    Count(ctx context.Context) (int64, error)
}
```

### 4.3 Service Interface (internal/domain/services/user_service.go)
```go
package services

import (
    "context"
    "github.com/naufalfazanadi/finance-manager-go/internal/application/dto"
)

type UserService interface {
    CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
    GetUser(ctx context.Context, id string) (*dto.UserResponse, error)
    GetUsers(ctx context.Context, limit, offset int) (*dto.UsersResponse, error)
    UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
    DeleteUser(ctx context.Context, id string) error
}
```

## Step 5: Application Layer

### 5.1 DTOs (internal/application/dto/user_dto.go)
```go
package dto

import (
    "time"
    "github.com/google/uuid"
)

// Request DTOs
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Age   int    `json:"age" validate:"min=1,max=120"`
}

type UpdateUserRequest struct {
    Name string `json:"name" validate:"omitempty,min=2,max=100"`
    Age  int    `json:"age" validate:"omitempty,min=1,max=120"`
}

// Response DTOs
type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UsersResponse struct {
    Users []UserResponse `json:"users"`
    Total int64          `json:"total"`
    Limit int            `json:"limit"`
    Offset int           `json:"offset"`
}

// Error Response
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

### 5.2 Use Cases (internal/application/usecases/user_usecase.go)
```go
package usecases

import (
    "context"
    "errors"
    "fmt"

    "github.com/naufalfazanadi/finance-manager-go/internal/application/dto"
    "github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
    "github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
    "github.com/naufalfazanadi/finance-manager-go/pkg/validator"

    "github.com/google/uuid"
)

type UserUseCase struct {
    userRepo  repositories.UserRepository
    validator *validator.Validator
}

func NewUserUseCase(userRepo repositories.UserRepository, validator *validator.Validator) *UserUseCase {
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

func (uc *UserUseCase) GetUsers(ctx context.Context, limit, offset int) (*dto.UsersResponse, error) {
    users, err := uc.userRepo.GetAll(ctx, limit, offset)
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

    return &dto.UsersResponse{
        Users:  userResponses,
        Total:  total,
        Limit:  limit,
        Offset: offset,
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
```

## Step 6: Infrastructure Layer

### 6.1 Database Connection (internal/infrastructure/database/postgres.go)
```go
package database

import (
    "fmt"
    "log"

    "github.com/naufalfazanadi/finance-manager-go/internal/config"
    "github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) *gorm.DB {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.DBName,
        cfg.Database.SSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate
    if err := db.AutoMigrate(&entities.User{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    log.Println("Database connected and migrated successfully")
    return db
}
```

### 6.2 Repository Implementation (internal/infrastructure/repositories/user_repository_impl.go)
```go
package repositories

import (
    "context"
    "errors"

    "github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
    "github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
    if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
        return err
    }
    return nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
    var user entities.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
    var user entities.User
    if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entities.User, error) {
    var users []*entities.User
    query := r.db.WithContext(ctx)

    if limit > 0 {
        query = query.Limit(limit)
    }
    if offset > 0 {
        query = query.Offset(offset)
    }

    if err := query.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
    if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
        return err
    }
    return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
    if err := r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error; err != nil {
        return err
    }
    return nil
}

func (r *userRepositoryImpl) Count(ctx context.Context) (int64, error) {
    var count int64
    if err := r.db.WithContext(ctx).Model(&entities.User{}).Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}
```

## Step 7: Interface Layer

### 7.1 HTTP Handlers (internal/interfaces/http/handlers/user_handler.go)
```go
package handlers

import (
    "strconv"

    "github.com/naufalfazanadi/finance-manager-go/internal/application/dto"
    "github.com/naufalfazanadi/finance-manager-go/internal/application/usecases"

    "github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    userUseCase *usecases.UserUseCase
}

func NewUserHandler(userUseCase *usecases.UserUseCase) *UserHandler {
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
        return c.Status(fiber.StatusBadRequest).JSON(&dto.ErrorResponse{
            Error:   "Bad Request",
            Message: "Invalid request body",
        })
    }

    user, err := h.userUseCase.CreateUser(c.Context(), &req)
    if err != nil {
        status := fiber.StatusInternalServerError
        if err.Error() == "user with this email already exists" {
            status = fiber.StatusConflict
        } else if err.Error()[:10] == "validation" {
            status = fiber.StatusBadRequest
        }

        return c.Status(status).JSON(&dto.ErrorResponse{
            Error:   "Creation Failed",
            Message: err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(user)
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
        return c.Status(fiber.StatusBadRequest).JSON(&dto.ErrorResponse{
            Error:   "Bad Request",
            Message: "User ID is required",
        })
    }

    user, err := h.userUseCase.GetUser(c.Context(), id)
    if err != nil {
        status := fiber.StatusInternalServerError
        if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
            status = fiber.StatusNotFound
        }

        return c.Status(status).JSON(&dto.ErrorResponse{
            Error:   "Not Found",
            Message: err.Error(),
        })
    }

    return c.JSON(user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.UsersResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    limit := 10
    offset := 0

    if l := c.Query("limit"); l != "" {
        if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
            limit = parsed
        }
    }

    if o := c.Query("offset"); o != "" {
        if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
            offset = parsed
        }
    }

    users, err := h.userUseCase.GetUsers(c.Context(), limit, offset)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(&dto.ErrorResponse{
            Error:   "Internal Server Error",
            Message: err.Error(),
        })
    }

    return c.JSON(users)
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
        return c.Status(fiber.StatusBadRequest).JSON(&dto.ErrorResponse{
            Error:   "Bad Request",
            Message: "User ID is required",
        })
    }

    var req dto.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(&dto.ErrorResponse{
            Error:   "Bad Request",
            Message: "Invalid request body",
        })
    }

    user, err := h.userUseCase.UpdateUser(c.Context(), id, &req)
    if err != nil {
        status := fiber.StatusInternalServerError
        if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
            status = fiber.StatusNotFound
        } else if err.Error()[:10] == "validation" {
            status = fiber.StatusBadRequest
        }

        return c.Status(status).JSON(&dto.ErrorResponse{
            Error:   "Update Failed",
            Message: err.Error(),
        })
    }

    return c.JSON(user)
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
        return c.Status(fiber.StatusBadRequest).JSON(&dto.ErrorResponse{
            Error:   "Bad Request",
            Message: "User ID is required",
        })
    }

    err := h.userUseCase.DeleteUser(c.Context(), id)
    if err != nil {
        status := fiber.StatusInternalServerError
        if err.Error() == "user not found" || err.Error() == "invalid user ID format" {
            status = fiber.StatusNotFound
        }

        return c.Status(status).JSON(&dto.ErrorResponse{
            Error:   "Delete Failed",
            Message: err.Error(),
        })
    }

    return c.SendStatus(fiber.StatusNoContent)
}
```

### 7.2 Routes (internal/interfaces/http/routes/routes.go)
```go
package routes

import (
    "github.com/naufalfazanadi/finance-manager-go/internal/interfaces/http/handlers"
    "github.com/naufalfazanadi/finance-manager-go/internal/interfaces/http/middleware"

    "github.com/gofiber/fiber/v2"
    fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

type Routes struct {
    userHandler *handlers.UserHandler
}

func NewRoutes(userHandler *handlers.UserHandler) *Routes {
    return &Routes{
        userHandler: userHandler,
    }
}

func (r *Routes) Setup() *fiber.App {
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler,
    })

    // Middleware
    app.Use(recover.New())
    app.Use(fiberLogger.New())
    app.Use(middleware.CORS())

    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":  "ok",
            "message": "Server is running",
        })
    })

    // API routes
    api := app.Group("/api/v1")

    // User routes
    users := api.Group("/users")
    users.Post("/", r.userHandler.CreateUser)
    users.Get("/", r.userHandler.GetUsers)
    users.Get("/:id", r.userHandler.GetUser)
    users.Put("/:id", r.userHandler.UpdateUser)
    users.Delete("/:id", r.userHandler.DeleteUser)

    return app
}
```

### 7.3 Middleware (internal/interfaces/http/middleware/middleware.go)
```go
package middleware

import (
    "github.com/naufalfazanadi/finance-manager-go/internal/application/dto"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
)

// ErrorHandler handles fiber errors
func ErrorHandler(c *fiber.Ctx, err error) error {
    code := fiber.StatusInternalServerError

    if e, ok := err.(*fiber.Error); ok {
        code = e.Code
    }

    return c.Status(code).JSON(&dto.ErrorResponse{
        Error:   "Error",
        Message: err.Error(),
    })
}

// CORS middleware
func CORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin,Content-Type,Accept,Authorization",
    })
}
```

## Step 8: Package Layer

### 8.1 Logger (pkg/logger/logger.go)
```go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(level string) error {
    config := zap.NewProductionConfig()

    switch level {
    case "debug":
        config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
    case "info":
        config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
    case "warn":
        config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
    case "error":
        config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
    default:
        config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
    }

    var err error
    Logger, err = config.Build()
    if err != nil {
        return err
    }

    return nil
}

func Info(msg string, fields ...zap.Field) {
    Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    Logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
    Logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
    Logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
    Logger.Fatal(msg, fields...)
}
```

### 8.2 Validator (pkg/validator/validator.go)
```go
package validator

import (
    "errors"
    "fmt"
    "strings"

    "github.com/go-playground/validator/v10"
)

type Validator struct {
    validate *validator.Validate
}

func New() *Validator {
    return &Validator{
        validate: validator.New(),
    }
}

func (v *Validator) Validate(i interface{}) error {
    if err := v.validate.Struct(i); err != nil {
        var validationErrors []string
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors = append(validationErrors, fmt.Sprintf(
                "field '%s' failed validation for tag '%s'",
                strings.ToLower(err.Field()),
                err.Tag(),
            ))
        }
        return errors.New(strings.Join(validationErrors, ", "))
    }
    return nil
}
```

## Step 9: Dependency Injection

### 9.1 Container (internal/config/container.go)
```go
package config

import (
    "github.com/naufalfazanadi/finance-manager-go/internal/application/usecases"
    "github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/database"
    "github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/repositories"
    "github.com/naufalfazanadi/finance-manager-go/internal/interfaces/http/handlers"
    "github.com/naufalfazanadi/finance-manager-go/pkg/validator"
)

type Container struct {
    UserHandler *handlers.UserHandler
}

func NewContainer(cfg *Config) *Container {
    // Database
    db := database.NewPostgresDB(cfg)

    // Validator
    validator := validator.New()

    // Repositories
    userRepo := repositories.NewUserRepository(db)

    // Use Cases
    userUseCase := usecases.NewUserUseCase(userRepo, validator)

    // Handlers
    userHandler := handlers.NewUserHandler(userUseCase)

    return &Container{
        UserHandler: userHandler,
    }
}
```

## Step 10: Main Application

### 10.1 Main File (cmd/server/main.go)
```go
package main

import (
    "log"

    "github.com/naufalfazanadi/finance-manager-go/internal/config"
    "github.com/naufalfazanadi/finance-manager-go/internal/interfaces/http/routes"
    "github.com/naufalfazanadi/finance-manager-go/pkg/logger"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize logger
    if err := logger.Init(cfg.App.LogLevel); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Logger.Sync()

    // Initialize container with dependencies
    container := config.NewContainer(cfg)

    // Setup routes
    routeHandler := routes.NewRoutes(container.UserHandler)
    app := routeHandler.Setup()

    // Start server
    serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
    logger.Info("Starting server on " + serverAddr)

    if err := app.Listen(":" + cfg.Server.Port); err != nil {
        logger.Fatal("Failed to start server: " + err.Error())
    }
}
```

## Step 11: Database Setup

### 11.1 Create PostgreSQL Database
```sql
-- Connect to PostgreSQL as superuser
CREATE DATABASE clean_api_db;
CREATE USER postgres WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE clean_api_db TO postgres;

-- Or use Docker
docker run --name postgres-clean-api \
  -e POSTGRES_DB=clean_api_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15
```

### 11.2 Database Migration Script (scripts/migrate.sh)
```bash
#!/bin/bash

# Database migration script
echo "Running database migrations..."

# You can add custom migration logic here
# For now, GORM AutoMigrate handles this

echo "Migrations completed!"
```

## Step 12: Testing the API

### 12.1 Run the Application
```bash
# Make sure PostgreSQL is running
# Update .env file with correct database credentials

# Run the application
go run cmd/server/main.go
```

### 12.2 Test API Endpoints

#### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "name": "John Doe",
    "age": 30
  }'
```

#### Get All Users
```bash
curl -X GET http://localhost:8080/api/v1/users?limit=10&offset=0
```

#### Get User by ID
```bash
curl -X GET http://localhost:8080/api/v1/users/{user-id}
```

#### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "age": 31
  }'
```

#### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id}
```

#### Health Check
```bash
curl -X GET http://localhost:8080/health
```

## Step 13: Docker Support (Optional)

### 13.1 Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

### 13.2 Docker Compose (docker-compose.yml)
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres-clean-api
    environment:
      POSTGRES_DB: clean_api_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - clean-api-network

  app:
    build: .
    container_name: clean-api-app
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: password
      DB_NAME: clean_api_db
      DB_SSLMODE: disable
      SERVER_PORT: 8080
      APP_ENV: production
      LOG_LEVEL: info
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    networks:
      - clean-api-network

volumes:
  postgres_data:

networks:
  clean-api-network:
    driver: bridge
```

## Step 14: Build and Run Commands

### 14.1 Build Script (scripts/build.sh)
```bash
#!/bin/bash

echo "Building the application..."

# Build for current platform
go build -o bin/server cmd/server/main.go

echo "Build completed! Binary is in bin/server"
```

### 14.2 Makefile
```makefile
.PHONY: build run test clean docker-up docker-down

# Build the application
build:
	@echo "Building the application..."
	@go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running the application..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/

# Start services with Docker Compose
docker-up:
	@echo "Starting services..."
	@docker-compose up -d

# Stop services
docker-down:
	@echo "Stopping services..."
	@docker-compose down

# Build and run with Docker
docker-build:
	@echo "Building Docker image..."
	@docker build -t clean-api .

# Database migration
migrate:
	@echo "Running migrations..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
```

## Important Notes

### 1. Clean Architecture Benefits
- **Separation of Concerns**: Each layer has a specific responsibility
- **Testability**: Easy to unit test business logic
- **Independence**: Business rules don't depend on frameworks
- **Flexibility**: Easy to change external dependencies

### 2. GORM Best Practices
- Use context for all database operations
- Enable soft deletes with `gorm.DeletedAt`
- Use proper indexes for performance
- Handle errors appropriately

### 3. Fiber Framework Features
- Fast HTTP router
- Built-in middleware support
- Easy JSON handling
- Excellent performance

### 4. Project Structure Benefits
- **Domain Layer**: Contains business entities and rules
- **Application Layer**: Orchestrates use cases
- **Infrastructure Layer**: Handles external dependencies
- **Interface Layer**: Manages external communication

### 5. Error Handling
- Use structured error responses
- Validate input data
- Handle database errors gracefully
- Return appropriate HTTP status codes

### 6. Security Considerations
- Validate all input data
- Use environment variables for sensitive data
- Implement proper error handling
- Add authentication/authorization as needed

### 7. Performance Tips
- Use database indexes
- Implement pagination
- Use connection pooling
- Add caching layer if needed

## Next Steps

1. **Add Authentication**: Implement JWT-based authentication
2. **Add Logging**: Enhance logging with structured logs
3. **Add Testing**: Write unit and integration tests
4. **Add Monitoring**: Implement health checks and metrics
5. **Add Documentation**: Generate API documentation with Swagger
6. **Add Caching**: Implement Redis for caching
7. **Add Rate Limiting**: Protect against abuse
8. **Add Validation**: Enhanced input validation

This tutorial provides a solid foundation for a production-ready Go backend with Clean Architecture principles.
