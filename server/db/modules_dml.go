package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/TLop503/LogCrunch/structs"
	"time"
)

// Module represents a parser module and its schema
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

// LoadModulesFromRegistry adds the contents of the metaparser reg
// to the database
func LoadModulesFromRegistry(db *sql.DB) error {
	for name, entry := range structs.MetaParserRegistry {
		//Marshal to json
		schemaJson, err := json.Marshal(entry.Schema)
		if err != nil {
			return fmt.Errorf("Failed to marshall schema for fodule %s: %w", name, schemaJson)
		}

		err = InsertModule(db, DBModule{Name: name, Schema: string(schemaJson)})
		if err != nil {
			return fmt.Errorf("Error inserting module to db: %w", err)
		}
	}
	return nil
}
