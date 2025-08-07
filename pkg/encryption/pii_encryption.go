package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

// EncryptionResult represents the result wrapper
type EncryptionResult struct {
	Data  interface{}
	Error error
}

// getSecretKey gets the encryption secret key from environment
func getSecretKey() string {
	return getEnv("ENCRYPTION_SECRET_KEY", "")
}

// getPepper gets the encryption pepper from environment
func getPepper() string {
	return getEnv("ENCRYPTION_PEPPER", "default-pepper-change-in-production")
}

// HashSHA256 creates a SHA256 hash of the text with pepper
func HashSHA256(text string) EncryptionResult {
	if text == "" {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("text required"),
		}
	}

	pepper := getPepper()
	hash := sha256.New()
	hash.Write([]byte(text + pepper))
	result := fmt.Sprintf("%x", hash.Sum(nil))

	return EncryptionResult{
		Data:  result,
		Error: nil,
	}
}

// EncryptAES128GCM encrypts text using AES-128-GCM
func EncryptAES128GCM(text string) EncryptionResult {
	if text == "" {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("text required"),
		}
	}

	secretKey := getSecretKey()
	if secretKey == "" {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("encryption secret key not configured"),
		}
	}

	// Decode the base64 secret key
	cipherKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to decode secret key: %w", err),
		}
	}

	// Create AES cipher
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to create cipher: %w", err),
		}
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to create GCM: %w", err),
		}
	}

	// Generate random IV (12 bytes for GCM)
	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to generate IV: %w", err),
		}
	}

	// Encrypt the text
	ciphertext := gcm.Seal(nil, iv, []byte(text), nil)

	// Combine IV + ciphertext + tag (tag is already included in Seal output)
	result := append(iv, ciphertext...)

	return EncryptionResult{
		Data:  result,
		Error: nil,
	}
}

// DecryptAES128GCM decrypts data using AES-128-GCM
func DecryptAES128GCM(data interface{}) EncryptionResult {
	if data == nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("input required"),
		}
	}

	secretKey := getSecretKey()
	if secretKey == "" {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("encryption secret key not configured"),
		}
	}

	var dataBytes []byte
	var err error

	// Handle different input types
	switch v := data.(type) {
	case []byte:
		dataBytes = v
	case string:
		// Assume string is base64 encoded
		dataBytes, err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			return EncryptionResult{
				Data:  nil,
				Error: fmt.Errorf("failed to decode base64 string: %w", err),
			}
		}
	default:
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("unsupported data type"),
		}
	}

	// Minimum size check: 12 bytes IV + 16 bytes tag = 28 bytes minimum
	if len(dataBytes) < 28 {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("invalid data length"),
		}
	}

	// Decode the base64 secret key
	cipherKey, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to decode secret key: %w", err),
		}
	}

	// Create AES cipher
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to create cipher: %w", err),
		}
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to create GCM: %w", err),
		}
	}

	// Extract IV (first 12 bytes) and ciphertext (remaining bytes)
	iv := dataBytes[:12]
	ciphertext := dataBytes[12:]

	// Decrypt
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return EncryptionResult{
			Data:  nil,
			Error: fmt.Errorf("failed to decrypt: %w", err),
		}
	}

	return EncryptionResult{
		Data:  string(plaintext),
		Error: nil,
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
