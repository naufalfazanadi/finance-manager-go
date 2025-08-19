package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	mathRand "math/rand"
	"time"
)

// generateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// EncryptResetToken encrypts a token with timestamp for forgot password functionality
func EncryptResetToken(randomString string) (string, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return "", fmt.Errorf("encryption secret key not configured")
	}

	// Decode the base64 secret key
	cipherKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret key: %w", err)
	}

	// Create payload: randomString.timestamp
	timestamp := time.Now().UnixMilli()
	payload := fmt.Sprintf("%s.%d", randomString, timestamp)

	// Create AES cipher
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random IV (12 bytes for GCM)
	iv := make([]byte, 12)
	if _, err := cryptoRand.Read(iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	// Encrypt the payload
	ciphertext := gcm.Seal(nil, iv, []byte(payload), nil)

	// Combine IV + ciphertext
	result := append(iv, ciphertext...)

	// Return base64 encoded string
	return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptResetToken decrypts a token and returns the random string and timestamp
func DecryptResetToken(encryptedToken string) (randomString string, timestamp int64, err error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return "", 0, fmt.Errorf("encryption secret key not configured")
	}

	// Decode the base64 secret key
	cipherKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode secret key: %w", err)
	}

	// Decode base64
	dataBytes, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode base64 token: %w", err)
	}

	// Minimum size check: 12 bytes IV + 16 bytes tag = 28 bytes minimum
	if len(dataBytes) < 28 {
		return "", 0, fmt.Errorf("invalid token length")
	}

	// Create AES cipher
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract IV (first 12 bytes) and ciphertext (remaining bytes)
	iv := dataBytes[:12]
	ciphertext := dataBytes[12:]

	// Decrypt
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Parse the payload (randomString.timestamp)
	payload := string(plaintext)
	parts := splitLast(payload, ".")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid token format")
	}

	// Parse timestamp
	timestamp, err = parseTimestamp(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid timestamp in token: %w", err)
	}

	return parts[0], timestamp, nil
}

// ValidateResetTokenExpiry checks if the token is still valid based on the expiry duration
func ValidateResetTokenExpiry(timestamp int64, expiryDurationMs int64) error {
	currentTime := time.Now().UnixMilli()
	if currentTime-timestamp >= expiryDurationMs {
		return fmt.Errorf("token has expired")
	}
	return nil
}

// CheckResetTokenCooldown checks if enough time has passed since the last token generation
func CheckResetTokenCooldown(timestamp int64, cooldownDurationMs int64) error {
	currentTime := time.Now().UnixMilli()
	if currentTime-timestamp <= cooldownDurationMs {
		remainingSeconds := int((cooldownDurationMs - (currentTime - timestamp)) / 1000)
		return fmt.Errorf("please wait %d seconds before requesting another password reset", remainingSeconds)
	}
	return nil
}

// splitLast splits a string by delimiter and returns the last two parts
func splitLast(s, delimiter string) []string {
	parts := make([]string, 0, 2)
	lastIndex := -1

	// Find the last occurrence of delimiter
	for i := len(s) - 1; i >= 0; i-- {
		if s[i:i+len(delimiter)] == delimiter {
			lastIndex = i
			break
		}
	}

	if lastIndex == -1 {
		return []string{s}
	}

	parts = append(parts, s[:lastIndex])
	parts = append(parts, s[lastIndex+len(delimiter):])

	return parts
}

// parseTimestamp parses a timestamp string to int64
func parseTimestamp(timestampStr string) (int64, error) {
	var timestamp int64
	_, err := fmt.Sscanf(timestampStr, "%d", &timestamp)
	return timestamp, err
}
