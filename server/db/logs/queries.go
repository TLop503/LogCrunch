package logs

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/TLop503/LogCrunch/structs"
	_ "modernc.org/sqlite"
)

// MostRecent50 returns the 50 most recent log entries
func MostRecent50_old(db *sql.DB) ([]structs.Log, error) {
	stmt := `
	SELECT timestamp, name, host, parsed, raw
	FROM logs
	ORDER BY timestamp DESC
	LIMIT 50
	`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []structs.Log
	for rows.Next() {
		var log structs.Log
		err = rows.Scan(&log.Timestamp, &log.Name, &log.Host, &log.Parsed, &log.Raw)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// MostRecent50 debugger
func MostRecent50(db *sql.DB) ([]structs.Log, error) {
	stmt := `
	SELECT timestamp, name, host, parsed, raw
	FROM logs
	ORDER BY timestamp DESC
	LIMIT 50
	`

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("MostRecent50: query failed: %w", err)
	}
	defer rows.Close()

	var logs []structs.Log
	rowNum := 0

	for rows.Next() {
		rowNum++

		var (
			log       structs.Log
			parsedRaw any
		)

		if err := rows.Scan(
			&log.Timestamp,
			&log.Name,
			&log.Host,
			&parsedRaw,
			&log.Raw,
		); err != nil {
			return nil, fmt.Errorf(
				"MostRecent50: scan failed on row %d: %w",
				rowNum, err,
			)
		}

		// Inspect what SQLite gave us
		switch v := parsedRaw.(type) {
		case []byte:
			if err := json.Unmarshal(v, &log.Parsed); err != nil {
				return nil, fmt.Errorf(
					"MostRecent50: JSON unmarshal failed on row %d (parsed is []byte): %w; raw=%q",
					rowNum, err, string(v),
				)
			}
		case string:
			if err := json.Unmarshal([]byte(v), &log.Parsed); err != nil {
				return nil, fmt.Errorf(
					"MostRecent50: JSON unmarshal failed on row %d (parsed is string): %w; raw=%q",
					rowNum, err, v,
				)
			}
		case nil:
			return nil, fmt.Errorf(
				"MostRecent50: parsed column is NULL on row %d",
				rowNum,
			)
		default:
			return nil, fmt.Errorf(
				"MostRecent50: unexpected type for parsed column on row %d: %T (%v)",
				rowNum, v, v,
			)
		}

		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("MostRecent50: row iteration failed: %w", err)
	}

	return logs, nil
}

// RunQuery executes a custom query and returns log entries
func RunQuery(db *sql.DB, query string) ([]structs.Log, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []structs.Log
	for rows.Next() {
		var log structs.Log
		err = rows.Scan(&log.Timestamp, &log.Name, &log.Host, &log.Parsed, &log.Raw)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}
