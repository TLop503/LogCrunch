package users_test

import (
	"os"
	"testing"
	"time"

	"github.com/TLop503/LogCrunch/server/db/users"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (string, func()) {
	tmpFile, err := os.CreateTemp("", "test_users_*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	return dbPath, func() {
		os.Remove(dbPath)
	}
}

func TestInitUserDB(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Verify users table exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		t.Error("users table does not exist")
	}

	// Verify sessions table exists
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='sessions'").Scan(&tableName)
	if err != nil {
		t.Error("sessions table does not exist")
	}

	// Verify indexes exist
	expectedIndexes := []string{
		"idx_sessions_user_id",
		"idx_sessions_expires_at",
		"idx_users_username",
	}

	for _, idx := range expectedIndexes {
		var indexName string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name=?", idx).Scan(&indexName)
		if err != nil {
			t.Errorf("Index %s does not exist", idx)
		}
	}
}

func TestCreateUser(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create a user
	id, err := users.CreateUser(db, "admin", "hashed_password_123", true)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected user ID 1, got %d", id)
	}

	// Verify user was created
	user, err := users.GetUserByUsername(db, "admin")
	if err != nil {
		t.Fatalf("GetUserByUsername failed: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", user.Username)
	}
	if user.PasswordHash != "hashed_password_123" {
		t.Errorf("Password hash mismatch")
	}
	if !user.CanCreateUsers {
		t.Error("Expected can_create_users to be true")
	}
	if !user.RequiresPasswordChange {
		t.Error("Expected requires_password_change to be true for new user")
	}
	if !user.IsActive {
		t.Error("Expected is_active to be true")
	}
}

func TestCreateUserDuplicateUsername(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create first user
	_, err = users.CreateUser(db, "admin", "password1", true)
	if err != nil {
		t.Fatalf("First CreateUser failed: %v", err)
	}

	// Attempt to create duplicate user
	_, err = users.CreateUser(db, "admin", "password2", false)
	if err == nil {
		t.Error("Expected error when creating duplicate username")
	}
}

func TestGetUserByID(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create a user
	id, err := users.CreateUser(db, "testuser", "hash123", false)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Get by ID
	user, err := users.GetUserByID(db, id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	// Get non-existent user
	user, err = users.GetUserByID(db, 9999)
	if err != nil {
		t.Fatalf("GetUserByID failed for non-existent: %v", err)
	}
	if user != nil {
		t.Error("Expected nil for non-existent user")
	}
}

func TestUpdatePassword(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create a user
	id, err := users.CreateUser(db, "admin", "old_hash", true)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Verify requires_password_change is true
	user, _ := users.GetUserByID(db, id)
	if !user.RequiresPasswordChange {
		t.Error("Expected requires_password_change to be true initially")
	}

	// Update password
	err = users.UpdatePassword(db, id, "new_hash")
	if err != nil {
		t.Fatalf("UpdatePassword failed: %v", err)
	}

	// Verify password was updated and requires_password_change is false
	user, err = users.GetUserByID(db, id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.PasswordHash != "new_hash" {
		t.Error("Password hash not updated")
	}
	if user.RequiresPasswordChange {
		t.Error("Expected requires_password_change to be false after update")
	}
}

func TestUpdateLastLogin(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create a user
	id, err := users.CreateUser(db, "admin", "hash", true)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Update last login
	err = users.UpdateLastLogin(db, id, "192.168.1.100")
	if err != nil {
		t.Fatalf("UpdateLastLogin failed: %v", err)
	}

	// Verify last login was updated
	user, err := users.GetUserByID(db, id)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.LastSeenIP == nil || *user.LastSeenIP != "192.168.1.100" {
		t.Error("Last seen IP not updated correctly")
	}
	if user.LastLogin == nil {
		t.Error("Last login timestamp not set")
	}
}

func TestSetUserActive(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create a user
	id, err := users.CreateUser(db, "admin", "hash", true)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Deactivate user
	err = users.SetUserActive(db, id, false)
	if err != nil {
		t.Fatalf("SetUserActive failed: %v", err)
	}

	// Verify user is deactivated
	user, _ := users.GetUserByID(db, id)
	if user.IsActive {
		t.Error("Expected user to be inactive")
	}

	// Reactivate user
	err = users.SetUserActive(db, id, true)
	if err != nil {
		t.Fatalf("SetUserActive failed: %v", err)
	}

	// Verify user is active
	user, _ = users.GetUserByID(db, id)
	if !user.IsActive {
		t.Error("Expected user to be active")
	}
}

func TestUserCount(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Initially zero
	count, err := users.UserCount(db)
	if err != nil {
		t.Fatalf("UserCount failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 users, got %d", count)
	}

	// Create users
	users.CreateUser(db, "user1", "hash1", false)
	users.CreateUser(db, "user2", "hash2", false)

	count, err = users.UserCount(db)
	if err != nil {
		t.Fatalf("UserCount failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 users, got %d", count)
	}
}

func TestIsInitialized(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Initially not initialized
	initialized, err := users.IsInitialized(db)
	if err != nil {
		t.Fatalf("IsInitialized failed: %v", err)
	}
	if initialized {
		t.Error("Expected not initialized")
	}

	// Create a user
	users.CreateUser(db, "admin", "hash", true)

	// Now initialized
	initialized, err = users.IsInitialized(db)
	if err != nil {
		t.Fatalf("IsInitialized failed: %v", err)
	}
	if !initialized {
		t.Error("Expected initialized")
	}
}

// Session tests

func TestCreateSession(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Create a user first
	userID, err := users.CreateUser(db, "admin", "hash", true)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Create a session
	session, err := users.CreateSession(db, userID, "192.168.1.100", time.Hour)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if session == nil {
		t.Fatal("Expected session, got nil")
	}
	if len(session.ID) != 64 { // 32 bytes hex encoded
		t.Errorf("Expected 64 char session ID, got %d", len(session.ID))
	}
	if session.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, session.UserID)
	}
	if session.IPAddress != "192.168.1.100" {
		t.Errorf("Expected IP '192.168.1.100', got '%s'", session.IPAddress)
	}
}

func TestGetSession(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	userID, _ := users.CreateUser(db, "admin", "hash", true)
	session, _ := users.CreateSession(db, userID, "192.168.1.100", time.Hour)

	// Get session
	retrieved, err := users.GetSession(db, session.ID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("Expected session, got nil")
	}
	if retrieved.ID != session.ID {
		t.Error("Session ID mismatch")
	}

	// Get non-existent session
	retrieved, err = users.GetSession(db, "nonexistent")
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if retrieved != nil {
		t.Error("Expected nil for non-existent session")
	}
}

func TestValidateSession(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	userID, _ := users.CreateUser(db, "admin", "hash", true)
	session, _ := users.CreateSession(db, userID, "192.168.1.100", time.Hour)

	// Valid session with matching IP
	valid, err := users.ValidateSession(db, session.ID, "192.168.1.100")
	if err != nil {
		t.Fatalf("ValidateSession failed: %v", err)
	}
	if valid == nil {
		t.Error("Expected valid session")
	}

	// Invalid IP
	valid, err = users.ValidateSession(db, session.ID, "192.168.1.200")
	if err != nil {
		t.Fatalf("ValidateSession failed: %v", err)
	}
	if valid != nil {
		t.Error("Expected nil for wrong IP")
	}

	// Non-existent session
	valid, err = users.ValidateSession(db, "nonexistent", "192.168.1.100")
	if err != nil {
		t.Fatalf("ValidateSession failed: %v", err)
	}
	if valid != nil {
		t.Error("Expected nil for non-existent session")
	}
}

func TestDeleteSession(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	userID, _ := users.CreateUser(db, "admin", "hash", true)
	session, _ := users.CreateSession(db, userID, "192.168.1.100", time.Hour)

	// Delete session
	err = users.DeleteSession(db, session.ID)
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	// Verify session is gone
	retrieved, _ := users.GetSession(db, session.ID)
	if retrieved != nil {
		t.Error("Expected session to be deleted")
	}
}

func TestDeleteAllUserSessions(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	userID, _ := users.CreateUser(db, "admin", "hash", true)

	// Create multiple sessions
	session1, _ := users.CreateSession(db, userID, "192.168.1.100", time.Hour)
	session2, _ := users.CreateSession(db, userID, "192.168.1.101", time.Hour)

	// Delete all sessions
	err = users.DeleteAllUserSessions(db, userID)
	if err != nil {
		t.Fatalf("DeleteAllUserSessions failed: %v", err)
	}

	// Verify both sessions are gone
	retrieved1, _ := users.GetSession(db, session1.ID)
	retrieved2, _ := users.GetSession(db, session2.ID)
	if retrieved1 != nil || retrieved2 != nil {
		t.Error("Expected all sessions to be deleted")
	}
}

func TestGetUserActiveSessions(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	userID, _ := users.CreateUser(db, "admin", "hash", true)

	// Create multiple sessions
	users.CreateSession(db, userID, "192.168.1.100", time.Hour)
	users.CreateSession(db, userID, "192.168.1.101", time.Hour)

	// Get active sessions
	sessions, err := users.GetUserActiveSessions(db, userID)
	if err != nil {
		t.Fatalf("GetUserActiveSessions failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("Expected 2 sessions, got %d", len(sessions))
	}
}

func TestForeignKeyConstraint(t *testing.T) {
	dbPath, cleanup := setupTestDB(t)
	defer cleanup()

	db, err := users.InitUserDB(dbPath)
	if err != nil {
		t.Fatalf("InitUserDB failed: %v", err)
	}
	defer db.Close()

	// Try to create session for non-existent user
	_, err = users.CreateSession(db, 9999, "192.168.1.100", time.Hour)
	if err == nil {
		t.Error("Expected error when creating session for non-existent user")
	}
}

func TestGenerateSessionID(t *testing.T) {
	// Generate multiple session IDs and verify uniqueness
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := users.GenerateSessionID()
		if err != nil {
			t.Fatalf("GenerateSessionID failed: %v", err)
		}
		if len(id) != 64 {
			t.Errorf("Expected 64 char ID, got %d", len(id))
		}
		if ids[id] {
			t.Error("Generated duplicate session ID")
		}
		ids[id] = true
	}
}
