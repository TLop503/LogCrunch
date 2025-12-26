package userauth

import (
	"crypto/rand"
	"fmt"
)

// generateDefaultPassword creates a random password for the initial admin user
func generateDefaultPassword() (string, error) {
	// For first-time setup, generate a random password
	passwordLen := 32
	bytes := make([]byte, passwordLen)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes), nil
}
