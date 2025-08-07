package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/pkg/config"
)

var (
	jwtSecret []byte
	jwtOnce   sync.Once
)

// initJWT initializes JWT secret (called once)
func initJWT() {
	jwtOnce.Do(func() {
		cfg := config.LoadConfig()
		secret := cfg.JWT.Secret
		// fmt.Printf("JWT initialized with secret length: %d\n", len(secret))
		jwtSecret = []byte(secret)
	})
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID uuid.UUID         `json:"user_id"`
	Email  string            `json:"email"`
	Name   string            `json:"name"`
	Role   entities.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the user
func GenerateToken(user *entities.User) (string, error) {
	initJWT() // Ensure JWT is initialized

	cfg := config.LoadConfig()

	now := time.Now()
	// Parse expiration time from config (default: 24h)
	expirationDuration, err := time.ParseDuration(cfg.JWT.ExpiresIn)
	if err != nil {
		// Fallback to 24 hours if parsing fails
		expirationDuration = 24 * time.Hour
		// fmt.Printf("Warning: Failed to parse JWT_EXPIRES_IN, using default 24h: %v\n", err)
	}
	expirationTime := now.Add(expirationDuration)

	// fmt.Printf("Token generation - Current time: %v (Unix: %d)\n", now, now.Unix())
	// fmt.Printf("Token generation - Expiration time: %v (Unix: %d)\n", expirationTime, expirationTime.Unix())

	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "finance-manager-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// fmt.Printf("Token generated successfully\n")
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	initJWT() // Ensure JWT is initialized

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		// fmt.Printf("JWT validation error: %v\n", err)
		return nil, err
	}

	if !token.Valid {
		// fmt.Printf("Token is not valid\n")
		return nil, errors.New("invalid token")
	}

	// fmt.Printf("Token validated successfully\n")
	return claims, nil
}

// ValidateTokenWithDB validates a JWT token and checks if user exists in database
// Returns updated claims with fresh data from database
func ValidateTokenWithDB(ctx context.Context, tokenString string, userRepo repositories.UserRepository) (*JWTClaims, error) {
	// First, validate the token signature and expiration
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Then, check if the user still exists in the database
	user, err := userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		// fmt.Printf("User validation failed: %v\n", err)
		return nil, errors.New("user not found or account has been deactivated")
	}

	// Update claims with fresh data from database (in case role or other info changed)
	claims.Email = user.Email
	claims.Name = user.Name
	claims.Role = user.Role

	// fmt.Printf("Token validated with database successfully for user: %s\n", user.Email)
	return claims, nil
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}
