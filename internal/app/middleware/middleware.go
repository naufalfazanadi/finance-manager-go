package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/auth"
)

// AuthMiddleware wraps authentication services
type AuthMiddleware struct{}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
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
		AllowOrigins: "*",
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
			return c.Status(fiber.StatusUnauthorized).JSON(&dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authorization header is required",
			})
		}

		// Extract token from header
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid authorization header format",
			})
		}

		// Validate token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&dto.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid or expired token",
			})
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
			return c.Status(fiber.StatusForbidden).JSON(&dto.ErrorResponse{
				Error:   "Forbidden",
				Message: "User role not found",
			})
		}

		currentRole := entities.UserRole(userRole.(string))

		// Check if user has any of the required roles
		for _, role := range roles {
			if currentRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(&dto.ErrorResponse{
			Error:   "Forbidden",
			Message: "Insufficient permissions",
		})
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
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next()
		}

		// Extract token from header
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Next()
		}

		// Validate token
		claims, err := auth.ValidateToken(token)
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
