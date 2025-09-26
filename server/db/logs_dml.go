package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/TLop503/LogCrunch/structs"
)

// InsertLog, allowing FKs to not necessarily exist yet
func InsertLog(db *sql.DB, l structs.Log) error {

	stmt := `
	INSERT INTO logs (name, path, host, timestamp, module, raw, parsed)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// If parsed is a struct/map, convert to JSON
	parsedJSON, err := json.Marshal(l.Parsed)
	if err != nil {
		return fmt.Errorf("failed to marshal parsed field: %w", err)
	}

	_, err = db.Exec(
		stmt,
		l.Name,
		l.Path,
		l.Host,
		l.Timestamp,
		l.Module,
		l.Raw,
		string(parsedJSON),
	)

	return err
}

// insert many logs at a time in batches for high-throughput
func InsertLogsBatch(db *sql.DB, logs []structs.Log) error {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure foreign keys are enabled and deferred
	if _, err := tx.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	if _, err := tx.Exec(`PRAGMA defer_foreign_keys = ON;`); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to defer foreign keys: %w", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO logs (name, path, host, timestamp, module, raw, parsed)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, l := range logs {
		parsedJSON, err := json.Marshal(l.Parsed)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to marshal parsed field: %w", err)
		}

		if _, err := stmt.Exec(
			l.Name,
			l.Path,
			l.Host,
			l.Timestamp,
			l.Module,
			l.Raw,
			string(parsedJSON),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert log: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
