package heartbeatlogs

import (
	"encoding/json"
	"time"

	"github.com/TLop503/heartbeat0/structs"
)

func GenerateLog(host string, hb structs.Heartbeat) (string, error) {
	now := time.Now().Unix()
	log := structs.Log{
		Host:      host,
		Timestamp: now,
		Type:      "{HB : Good}",
		Payload:   hb,
	}

	// Marshal the log into JSON format
	logJSON, jsonErr := json.Marshal(log)
	if jsonErr != nil {
		return "", jsonErr
	}

	// Convert JSON bytes to a string and return
	return string(logJSON), nil
}
