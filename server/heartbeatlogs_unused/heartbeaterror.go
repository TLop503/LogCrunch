package heartbeatlogs_unused

import (
	"encoding/json"
	"time"

	"github.com/TLop503/LogCrunch/structs"
)

func GenerateSeqErrorLog(host string, exp int, recv int) (string, error) {
	//get current time:
	now := time.Now().Unix()

	seqErr := structs.HbSeqErr{
		ExpSeq:  exp,
		RecvSeq: recv,
	}

	log := structs.Log{
		Host:      host,
		Timestamp: now,
		Type:      "{HB : Seq_Err}",
		Payload:   seqErr,
	}

	// Marshal the log into JSON format
	logJSON, jsonErr := json.Marshal(log)
	if jsonErr != nil {
		return "", jsonErr
	}

	// Convert JSON bytes to a string and return
	return string(logJSON), nil
}
