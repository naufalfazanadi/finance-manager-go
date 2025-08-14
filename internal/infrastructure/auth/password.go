package auth

import (
	"regexp"

	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a plain text password with a hashed password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength validates password requirements:
// - At least 1 uppercase letter
// - At least 1 number
// - At least 1 special character/symbol
func ValidatePasswordStrength(password string) error {
	if password == "" {
		return helpers.NewBadRequestError("password is required", "")
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return helpers.NewBadRequestError("password must contain at least one uppercase letter", "")
	}

	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return helpers.NewBadRequestError("password must contain at least one number", "")
	}

	// Check for at least one special character/symbol
	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)
	if !hasSymbol {
		return helpers.NewBadRequestError("password must contain at least one special character", "")
	}

	return nil
}
