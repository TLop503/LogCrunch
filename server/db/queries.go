package db

import (
	"database/sql"
	"github.com/TLop503/LogCrunch/structs"
)

func MostRecent50(db *sql.DB) ([]structs.Log, error) {
	stmt := `
	SELECT timestamp, name, host, parsed
	FROM logs
	ORDER BY timestamp DESC
	LIMIT 50
	`

	rows, err := db.Query(stmt)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var logs []structs.Log
	for rows.Next() {
		var log structs.Log
		err = rows.Scan(&log.Timestamp, &log.Name, &log.Host, &log.Parsed)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

func RunQuery(db *sql.DB, query string) ([]structs.Log, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []structs.Log
	for rows.Next() {
		var log structs.Log
		err = rows.Scan(&log.Timestamp, &log.Name, &log.Host, &log.Parsed)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}
