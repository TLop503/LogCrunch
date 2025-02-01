package heartbeat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TLop503/heartbeat0/structs"
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
func Heartbeat(writer *bufio.Writer, hostname string) {
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

		// Send the JSON data with a newline
		_, err = writer.WriteString(string(jsonData) + "\n")
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			break
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing data:", err)
			break
		}

		time.Sleep(60 * time.Second) // Send a heartbeat every minute. TODO: make this easy to configure.
	}
}
