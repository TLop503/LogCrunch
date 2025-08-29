package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const createModulesTable = `
CREATE TABLE IF NOT EXISTS modules (
    module_id   INTEGER PRIMARY KEY AUTOINCREMENT,
    module      TEXT NOT NULL UNIQUE,
    schema_json JSON NOT NULL,
    created_at  INTEGER NOT NULL DEFAULT (strftime('%s','now'))
);`

const createLogsTable = `
CREATE TABLE IF NOT EXISTS logs (
    log_id     INTEGER PRIMARY KEY AUTOINCREMENT,
    name 	   TEXT NOT NULL,
    path       TEXT NOT NULL,
    host       TEXT NOT NULL,
    timestamp  INTEGER NOT NULL,
    module     TEXT NOT NULL,
    raw        TEXT NOT NULL,
    parsed     JSON NOT NULL,
    FOREIGN KEY (module) REFERENCES modules(module)
);`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_logs_type ON logs(name);
CREATE INDEX IF NOT EXISTS idx_logs_host ON logs(host);
`

// InitDB initializes the SQLite database with tables and indexes.
// dbPath is the path to the .sqlite file.
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	statements := []string{
		createModulesTable,
		createLogsTable,
		createIndexes,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to execute statement: %w\nstmt: %s", err, stmt)
		}
	}

	return db, nil
}
