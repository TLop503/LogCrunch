package logs

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/TLop503/LogCrunch/structs"
)

// InsertLog inserts a single log entry, allowing FKs to not necessarily exist yet
func InsertLog(db *sql.DB, l structs.Log) error {
	parsedJSON, err := json.Marshal(l.Parsed)
	if err != nil {
		return fmt.Errorf("failed to marshal parsed field: %w", err)
	}

	// Ensure module exists FIRST
	if err := ensureModuleExists(db, l.Module, []byte(`{}`)); err != nil {
		return fmt.Errorf("failed to ensure module exists: %w", err)
	}

	_, err = db.Exec(`
		INSERT INTO logs (name, path, host, timestamp, module, raw, parsed)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
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

// InsertLogsBatch inserts many logs at a time in batches for high-throughput
func InsertLogsBatch(db *sql.DB, logs []structs.Log) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Always enable FKs per connection
	if _, err := tx.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	// 1. Collect unique modules
	modules := make(map[string]struct{})
	for _, l := range logs {
		modules[l.Module] = struct{}{}
	}

	// 2. Ensure all modules exist
	for module := range modules {
		if err := ensureModuleExists(tx, module, []byte(`{}`)); err != nil {
			return fmt.Errorf("failed to ensure module %q exists: %w", module, err)
		}
	}

	// 3. Insert logs
	stmt, err := tx.Prepare(`
		INSERT INTO logs (name, path, host, timestamp, module, raw, parsed)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range logs {
		parsedJSON, err := json.Marshal(l.Parsed)
		if err != nil {
			return err
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
			return err
		}
	}

	return tx.Commit()
}

// execer middleman to distinguish batch/solo workers
type execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// ensureModuleExists Upserts modules if an agent declares a new method of parsing
func ensureModuleExists(exec execer, module string, schemaJSON []byte) error {
	_, err := exec.Exec(`
		INSERT INTO modules (module, schema_json)
		VALUES (?, ?)
		ON CONFLICT(module) DO NOTHING;
	`, module, schemaJSON)

	return err
}
