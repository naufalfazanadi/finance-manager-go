package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/auth"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
)

// AuthMiddleware wraps authentication services
type AuthMiddleware struct {
	userRepo repositories.UserRepository
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
	}
}

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
		AllowOrigins: config.LoadConfig().CORS.AllowOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	})
}

// JWTAuth middleware to protect routes with JWT authentication
func (am *AuthMiddleware) JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helpers.HandleError(c, helpers.NewUnauthorizedError("Unauthorized", "Authorization header is required"), "Authorization header is required")
		}

		// Extract token from header with improved parsing
		var token string

		// Handle different Authorization header formats
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		} else if strings.HasPrefix(authHeader, "bearer ") {
			token = authHeader[7:]
		} else {
			// If no Bearer prefix, assume the entire header is the token
			token = strings.TrimSpace(authHeader)
		}

		// Validate that we have a token
		if token == "" {
			return helpers.HandleError(c, helpers.NewUnauthorizedError("Unauthorized", "Empty token provided"), "Empty token provided")
		}

		// Validate token with database check
		claims, err := auth.ValidateTokenWithDB(c.Context(), token, am.userRepo)
		if err != nil {
			return helpers.HandleError(c, helpers.NewUnauthorizedError("Unauthorized", "Invalid or expired token: "+err.Error()), "Invalid or expired token")
		}

		// Set user information in context
		c.Locals("userID", claims.UserID.String())
		c.Locals("userEmail", claims.Email)
		c.Locals("userName", claims.Name)
		c.Locals("userRole", string(claims.Role))

		return c.Next()
	}
}

// RequireRole middleware to check if user has required role
func RequireRole(roles ...entities.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return helpers.HandleError(c, helpers.NewForbiddenError("Forbidden", "User role not found"), "User role not found")
		}

		currentRole := entities.UserRole(userRole.(string))

		// Check if user has any of the required roles
		for _, role := range roles {
			if currentRole == role {
				return c.Next()
			}
		}

		return helpers.HandleError(c, helpers.NewForbiddenError("Forbidden", "Insufficient permissions"), "Insufficient permissions")
	}
}

// RequireAdmin middleware to check if user is admin
func RequireAdmin() fiber.Handler {
	return RequireRole(entities.UserRoleAdmin)
}

// RequireUser middleware to check if user is user
func RequireUser() fiber.Handler {
	return RequireRole(entities.UserRoleUser)
}

// OptionalJWTAuth middleware that doesn't require authentication but sets user info if token is provided
func (am *AuthMiddleware) OptionalJWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Extract token from header with improved parsing
		var token string

		// Handle different Authorization header formats
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		} else if strings.HasPrefix(authHeader, "bearer ") {
			token = authHeader[7:]
		} else {
			// If no Bearer prefix, assume the entire header is the token
			token = authHeader
		}

		// If no valid token, continue without authentication
		if token == "" {
			return c.Next()
		}

		// Validate token with database check
		claims, err := auth.ValidateTokenWithDB(c.Context(), token, am.userRepo)
		if err != nil {
			return c.Next()
		}

		// Set user information in context
		c.Locals("userID", claims.UserID.String())
		c.Locals("userEmail", claims.Email)
		c.Locals("userName", claims.Name)
		c.Locals("userRole", string(claims.Role))

		return c.Next()
	}
}
