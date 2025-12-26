package logs

import (
	"database/sql"
	"fmt"

	"github.com/TLop503/LogCrunch/server/db/core"
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

const enableForeignKeys = `PRAGMA foreign_keys = ON;`
const deferForeignKeys = `PRAGMA defer_foreign_keys = ON;`

// logStatements contains all DDL statements needed for the logs database
var logStatements = []string{
	createModulesTable,
	createLogsTable,
	createIndexes,
	enableForeignKeys,
	deferForeignKeys,
}

// InitLogDB initializes the logs SQLite database with tables and indexes.
// dbPath is the path to the .sqlite file.
func InitLogDB(dbPath string) (*sql.DB, error) {
	db, err := core.InitDB(dbPath, logStatements)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize log database: %w", err)
	}
	return db, nil
}
