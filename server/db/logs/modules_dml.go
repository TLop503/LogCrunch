package logs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TLop503/LogCrunch/structs"
)

// DBModule represents a parser module and its schema
type DBModule struct {
	Name   string
	Schema string // JSON string
}

// InsertModule inserts a module into the database, replacing it if it already exists
func InsertModule(db *sql.DB, m DBModule) error {
	stmt := `
	INSERT INTO modules (module, schema_json, created_at)
	VALUES (?, ?, ?)
	ON CONFLICT(module) DO UPDATE SET
		schema_json=excluded.schema_json,
		created_at=excluded.created_at;
	`

	_, err := db.Exec(stmt, m.Name, m.Schema, time.Now().Unix())
	return err
}

// loadModulesFromRegistry adds the contents of the metaparser registry
// to the database
func loadModulesFromRegistry(db *sql.DB) error {
	for name, entry := range structs.MetaParserRegistry {
		// Marshal to json
		schemaJson, err := json.Marshal(entry.Schema)
		if err != nil {
			return fmt.Errorf("failed to marshal schema for module %s: %s", name, schemaJson)
		}

		err = InsertModule(db, DBModule{Name: name, Schema: string(schemaJson)})
		if err != nil {
			return fmt.Errorf("error inserting module to db: %w", err)
		}
	}

	err := printAllModules(db)
	if err != nil {
		return fmt.Errorf("error printing all modules: %w", err)
	}
	return nil
}

// printAllModules queries the modules table and prints each module name for debugging
func printAllModules(db *sql.DB) error {
	rows, err := db.Query(`SELECT module FROM modules`)
	if err != nil {
		return fmt.Errorf("failed to query modules table: %w", err)
	}
	defer rows.Close()

	fmt.Println("Modules in database:")
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan module name: %w", err)
		}
		fmt.Println("-", name)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	return nil
}
