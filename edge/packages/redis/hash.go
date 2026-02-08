package redis

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashValue returns a SHA-256 hex digest for a string.
func HashValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
