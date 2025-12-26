package userauth

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/TLop503/LogCrunch/server/db/users"
	"golang.org/x/crypto/bcrypt"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Check for users, and create first if none exist and dotfile does not exist
func FirstTimeSetupCheck(userDB *sql.DB, dotfilePath string) error {
	userCount, err := users.UserCount(userDB)
	if err != nil {
		return fmt.Errorf("failed to check user count: %w", err)
	}

	if userCount == 0 {
		// No users in DB, trigger first-time setup
		if _, err := os.Stat(dotfilePath); errors.Is(err, os.ErrNotExist) {
			// Dotfile does not exist, proceed with setup
			return setup(userDB, dotfilePath)
		}
		return fmt.Errorf("No users exist, but dotfile suggests setup already happened! Something evil may be afoot.")
	}
	//TODO: check there exists at least one *active* user
	return nil
}

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

// setup creates initial admin user and dotfile
func setup(userDB *sql.DB, dotfilePath string) error {
	// create default admin user password
	defaultPassword, err := generateDefaultPassword()
	if err != nil {
		return fmt.Errorf("failed to generate default password: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	passwordHash := string(hash)

	// Create the initial admin user in DB
	_, err = users.CreateUser(userDB, "admin", passwordHash, true)
	if err != nil {
		return fmt.Errorf("failed to create initial admin user: %w", err)
	}

	// Only create the dotfile if user creation succeeded
	f, err := os.Create(dotfilePath)
	if err != nil {
		return fmt.Errorf("failed to create dotfile: %w", err)
	}
	f.Close()

	fmt.Printf("SETUP: User {admin} created with password {%s}. UPDATE THIS IMMEDIATELY!\n", defaultPassword)
	return nil
}
