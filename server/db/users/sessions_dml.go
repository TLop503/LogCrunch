package users

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// Session represents an active user session
type Session struct {
	ID        string
	UserID    int64
	CreatedAt time.Time
	ExpiresAt time.Time
	IPAddress string
}

// GenerateSessionID creates a cryptographically secure random session ID
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate session id: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession creates a new session for a user
func CreateSession(db *sql.DB, userID int64, ipAddress string, duration time.Duration) (*Session, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().Unix()
	expiresAt := time.Now().Add(duration).Unix()

	stmt := `
	INSERT INTO sessions (id, user_id, created_at, expires_at, ip_address)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.Exec(stmt, sessionID, userID, createdAt, expiresAt, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
		IPAddress: ipAddress,
	}, nil
}

// GetSession retrieves a session by ID, returns nil if not found or expired
func GetSession(db *sql.DB, sessionID string) (*Session, error) {
	stmt := `
	SELECT id, user_id, created_at, expires_at, ip_address
	FROM sessions
	WHERE id = ? AND expires_at > ?
	`

	var session Session
	var createdAt, expiresAt int64
	err := db.QueryRow(stmt, sessionID, time.Now().Unix()).Scan(
		&session.ID,
		&session.UserID,
		&createdAt,
		&expiresAt,
		&session.IPAddress,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session.CreatedAt = time.Unix(createdAt, 0)
	session.ExpiresAt = time.Unix(expiresAt, 0)

	return &session, nil
}

// ValidateSession checks if a session is valid for the given IP address
// Returns the session if valid, nil if invalid/expired/wrong IP
func ValidateSession(db *sql.DB, sessionID, ipAddress string) (*Session, error) {
	session, err := GetSession(db, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	// Verify IP address matches
	if session.IPAddress != ipAddress {
		return nil, nil
	}

	return session, nil
}

// DeleteSession removes a session (logout)
func DeleteSession(db *sql.DB, sessionID string) error {
	stmt := `DELETE FROM sessions WHERE id = ?`

	_, err := db.Exec(stmt, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// DeleteAllUserSessions removes all sessions for a user (force logout everywhere)
func DeleteAllUserSessions(db *sql.DB, userID int64) error {
	stmt := `DELETE FROM sessions WHERE user_id = ?`

	_, err := db.Exec(stmt, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

// CleanExpiredSessions removes all expired sessions from the database
func CleanExpiredSessions(db *sql.DB) (int64, error) {
	stmt := `DELETE FROM sessions WHERE expires_at < ?`

	result, err := db.Exec(stmt, time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("failed to clean expired sessions: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

// ExtendSession extends a session's expiration time
func ExtendSession(db *sql.DB, sessionID string, duration time.Duration) error {
	newExpiry := time.Now().Add(duration).Unix()

	stmt := `
	UPDATE sessions 
	SET expires_at = ?
	WHERE id = ? AND expires_at > ?
	`

	result, err := db.Exec(stmt, newExpiry, sessionID, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to extend session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("session not found or expired")
	}

	return nil
}

// GetUserActiveSessions returns all active sessions for a user
func GetUserActiveSessions(db *sql.DB, userID int64) ([]Session, error) {
	stmt := `
	SELECT id, user_id, created_at, expires_at, ip_address
	FROM sessions
	WHERE user_id = ? AND expires_at > ?
	ORDER BY created_at DESC
	`

	rows, err := db.Query(stmt, userID, time.Now().Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		var createdAt, expiresAt int64
		if err := rows.Scan(&s.ID, &s.UserID, &createdAt, &expiresAt, &s.IPAddress); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		s.CreatedAt = time.Unix(createdAt, 0)
		s.ExpiresAt = time.Unix(expiresAt, 0)
		sessions = append(sessions, s)
	}

	return sessions, rows.Err()
}
