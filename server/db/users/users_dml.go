package users

import (
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the database
type User struct {
	ID                     int64
	Username               string
	PasswordHash           string
	CanCreateUsers         bool
	CreatedAt              time.Time
	LastLogin              *time.Time
	LastSeenIP             *string
	IsActive               bool
	RequiresPasswordChange bool
}

// CreateUser inserts a new user into the database.
// New users are created with requires_password_change = true by default.
func CreateUser(db *sql.DB, username, passwordHash string, canCreateUsers bool) (int64, error) {
	stmt := `
	INSERT INTO users (username, password_hash, can_create_users, requires_password_change)
	VALUES (?, ?, ?, 1)
	`

	result, err := db.Exec(stmt, username, passwordHash, canCreateUsers)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	stmt := `
	SELECT id, username, password_hash, can_create_users, created_at, 
	       last_login, last_seen_ip, is_active, requires_password_change
	FROM users
	WHERE username = ?
	`

	var user User
	var lastLogin, lastSeenIP sql.NullString

	err := db.QueryRow(stmt, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CanCreateUsers,
		&user.CreatedAt,
		&lastLogin,
		&lastSeenIP,
		&user.IsActive,
		&user.RequiresPasswordChange,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		t, _ := time.Parse(time.RFC3339, lastLogin.String)
		user.LastLogin = &t
	}
	if lastSeenIP.Valid {
		user.LastSeenIP = &lastSeenIP.String
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sql.DB, id int64) (*User, error) {
	stmt := `
	SELECT id, username, password_hash, can_create_users, created_at, 
	       last_login, last_seen_ip, is_active, requires_password_change
	FROM users
	WHERE id = ?
	`

	var user User
	var lastLogin, lastSeenIP sql.NullString

	err := db.QueryRow(stmt, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CanCreateUsers,
		&user.CreatedAt,
		&lastLogin,
		&lastSeenIP,
		&user.IsActive,
		&user.RequiresPasswordChange,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if lastLogin.Valid {
		t, _ := time.Parse(time.RFC3339, lastLogin.String)
		user.LastLogin = &t
	}
	if lastSeenIP.Valid {
		user.LastSeenIP = &lastSeenIP.String
	}

	return &user, nil
}

// UpdatePassword updates a user's password hash and clears requires_password_change
func UpdatePassword(db *sql.DB, userID int64, newPasswordHash string) error {
	stmt := `
	UPDATE users 
	SET password_hash = ?, requires_password_change = 0
	WHERE id = ?
	`

	result, err := db.Exec(stmt, newPasswordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateLastLogin updates the last_login timestamp and last_seen_ip for a user
func UpdateLastLogin(db *sql.DB, userID int64, ip string) error {
	stmt := `
	UPDATE users 
	SET last_login = CURRENT_TIMESTAMP, last_seen_ip = ?
	WHERE id = ?
	`

	_, err := db.Exec(stmt, ip, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// SetUserActive enables or disables a user account
func SetUserActive(db *sql.DB, userID int64, isActive bool) error {
	stmt := `
	UPDATE users 
	SET is_active = ?
	WHERE id = ?
	`

	result, err := db.Exec(stmt, isActive, userID)
	if err != nil {
		return fmt.Errorf("failed to set user active status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UserCount returns the total number of users in the database
func UserCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// IsInitialized checks if the database has been initialized with at least one user
func IsInitialized(db *sql.DB) (bool, error) {
	count, err := UserCount(db)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
