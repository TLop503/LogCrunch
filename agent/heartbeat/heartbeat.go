package heartbeat

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/TLop503/LogCrunch/structs"
)

func Heartbeat(logChan chan<- string, hostname string) {
	seq := 0
	for { //forever

		// write log payload as pure json
		var seqAsJson interface{}
		rawHb := `{"seq":` + strconv.Itoa(seq) + `}`
		err := json.Unmarshal([]byte(rawHb), &seqAsJson)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
		}

		hb := structs.Log{
			Host:      hostname,
			Timestamp: time.Now().Unix(),
			Type:      "Heartbeat",
			Parsed:    nil,
			Raw:       seqAsJson,
		}

		//convert our struct to JSON
		jsonData, err := json.Marshal(hb)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			break
		}

		// Send the JSON data to the writer via chan
		logChan <- string(jsonData)

		time.Sleep(60 * time.Second) // Send a heartbeat every minute. TODO: make this easy to configure.
		seq++
	}
}
