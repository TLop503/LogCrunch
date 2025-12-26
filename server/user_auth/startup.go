package userauth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/TLop503/LogCrunch/server/filehandler"
	"log"
	"os"

	"github.com/TLop503/LogCrunch/server/db/users"
	"golang.org/x/crypto/bcrypt"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Check for users, and create first if none exist and dotfile does not exist
func FirstTimeSetupCheck(userDBPath, dotfilePath string) (*sql.DB, error) {
	userDB, err := users.InitUserDB(userDBPath)
	if err != nil {
		log.Fatalf("Error initializing user DB connection: %v", err)
	}

	userCount, err := users.UserCount(userDB)
	if err != nil {
		return nil, fmt.Errorf("failed to check user count: %w", err)
	}

	if userCount == 0 {
		// No users in DB, trigger first-time setup
		if _, err := os.Stat(dotfilePath); errors.Is(err, os.ErrNotExist) {
			// Dotfile does not exist, proceed with setup
			log.Println("Conducting first-time user db setup...")
			err = setup(userDB, dotfilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize user database: %w", err)
			}
			log.Println("First time user db setup successful!")
		}
		return nil, fmt.Errorf("No users exist, but dotfile suggests setup already happened! Something evil may be afoot.")
	}
	// happy case - successfully initialized user DB, or else was already setup
	//TODO: check there exists at least one *active* user
	return userDB, nil
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
	err = filehandler.Create_if_needed(dotfilePath, 0o755, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create dotfile: %w", err)
	}

	fmt.Printf("SETUP: User {admin} created with password {%s}. UPDATE THIS IMMEDIATELY!\n", defaultPassword)
	return nil
}
