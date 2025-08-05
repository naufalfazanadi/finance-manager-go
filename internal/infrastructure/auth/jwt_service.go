package auth

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
)

var (
	jwtSecret []byte
	jwtOnce   sync.Once
)

// initJWT initializes JWT secret (called once)
func initJWT() {
	jwtOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-secret-key-change-in-production"
		}
		fmt.Printf("JWT initialized with secret length: %d\n", len(secret))
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

	now := time.Now()
	expirationTime := now.Add(24 * time.Hour)

	fmt.Printf("Token generation - Current time: %v (Unix: %d)\n", now, now.Unix())
	fmt.Printf("Token generation - Expiration time: %v (Unix: %d)\n", expirationTime, expirationTime.Unix())

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

	fmt.Printf("Token generated successfully\n")
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
		fmt.Printf("JWT validation error: %v\n", err)
		return nil, err
	}

	if !token.Valid {
		fmt.Printf("Token is not valid\n")
		return nil, errors.New("invalid token")
	}

	fmt.Printf("Token validated successfully\n")
	return claims, nil
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return authHeader[7:], nil
}
