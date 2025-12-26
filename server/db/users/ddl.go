package users

import (
	"database/sql"
	"fmt"

	"github.com/TLop503/LogCrunch/server/db/core"
)

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    can_create_users BOOLEAN NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    last_seen_ip TEXT,
    is_active BOOLEAN DEFAULT 1,
    requires_password_change BOOLEAN DEFAULT 0
);`

const createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    ip_address TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
`

const enableForeignKeys = `PRAGMA foreign_keys = ON;`

// userStatements contains all DDL statements needed for the users database
var userStatements = []string{
	enableForeignKeys,
	createUsersTable,
	createSessionsTable,
	createIndexes,
}

// InitUserDB initializes the users SQLite database with tables and indexes.
// dbPath is the path to the .sqlite file.
func InitUserDB(dbPath string) (*sql.DB, error) {
	db, err := core.InitDB(dbPath, userStatements)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user database: %w", err)
	}
	return db, nil
}
