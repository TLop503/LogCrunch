package heartbeat

import (
	"strconv"
	"time"

	"github.com/TLop503/LogCrunch/structs"
)

// Heartbeat creates and transmits an "I'm alive" log every minute
func Heartbeat(logChan chan<- structs.Log, hostname string) {
	seq := 0
	for {
		// Create raw JSON as a map
		seqAsJSON := map[string]int{"seq": seq}

		// Create the log struct
		hb := structs.Log{
			Host:      hostname,
			Timestamp: time.Now().Unix(),
			Type:      "Heartbeat",
			Path:      "self",
			Parsed:    seqAsJSON,
			Raw:       strconv.Itoa(seq),
		}

		// Send the structured log over the channel
		logChan <- hb

		time.Sleep(60 * time.Second)
		seq++
	}
}
