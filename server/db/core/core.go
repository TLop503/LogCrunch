package core

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// InitDB initializes a SQLite database at the given path and executes
// the provided DDL statements. This is a generic function that can be
// used by any database package to set up their specific schema.
func InitDB(dbPath string, statements []string) (*sql.DB, error) {
	/*
		err := filehandler.Create_if_needed(dbPath, 0o755, 0o644)
		if err != nil {
			return nil, err
		}
	*/

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to execute statement: %w\nstmt: %s", err, stmt)
		}
	}

	return db, nil
}
