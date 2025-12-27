package logs_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/TLop503/LogCrunch/server/db/logs"
	_ "modernc.org/sqlite"
)

func TestInitLogDB(t *testing.T) {
	// Create a temporary file for the SQLite database
	tmpFile, err := os.CreateTemp("", "test_logcrunch_*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	// Initialize the database
	sqlDB, _, err := logs.InitLogDB(dbPath)
	if err != nil {
		t.Fatalf("InitLogDB failed: %v", err)
	}
	defer sqlDB.Close()

	// Verify that the modules table exists
	if !tableExists(sqlDB, "modules") {
		t.Error("modules table does not exist")
	}

	// Verify that the logs table exists
	if !tableExists(sqlDB, "logs") {
		t.Error("logs table does not exist")
	}

	// Verify that indexes exist
	expectedIndexes := []string{
		"idx_logs_timestamp",
		"idx_logs_type",
		"idx_logs_host",
	}

	for _, idx := range expectedIndexes {
		if !indexExists(sqlDB, idx) {
			t.Errorf("Index %s does not exist", idx)
		}
	}
}

// tableExists checks if a table exists in the SQLite database
func tableExists(db *sql.DB, tableName string) bool {
	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&name)
	return err == nil && name == tableName
}

// indexExists checks if an index exists in the SQLite database
func indexExists(db *sql.DB, indexName string) bool {
	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&name)
	return err == nil && name == indexName
}
