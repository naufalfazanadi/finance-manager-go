package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken creates a hash of the JWT token for use as a cache key
// This ensures sensitive token data is not stored directly as keys in Redis
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
