package heartbeat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TLop503/LogCrunch/structs"
)

// format the actual log payload into a struct
func makeHeartbeat(typ string, seq int, hostname string) structs.Heartbeat {
	return structs.Heartbeat{
		Type:      typ,
		Host:      hostname,
		Timestamp: time.Now().Unix(),
		Seq:       seq,
	}
}

// create hb via makeHeartbeat, then write to writer (which is sent over wire).
func Heartbeat(logChan chan<- string, hostname string) {
	seq := 0
	for { //forever

		hb := makeHeartbeat("proof_of_life", seq, hostname)
		seq++

		//convert our struct to JSON
		jsonData, err := json.Marshal(hb)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			break
		}

		// Send the JSON data to the writer via chan
		logChan <- string(jsonData)

		time.Sleep(60 * time.Second) // Send a heartbeat every minute. TODO: make this easy to configure.
	}
}
