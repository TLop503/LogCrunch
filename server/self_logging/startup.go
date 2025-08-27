package self_logging

import (
	"encoding/json"
	"fmt"
	"github.com/TLop503/LogCrunch/structs"
	"log"
	"os"
	"time"
)

// Generates a log for the firehose of the server start
func CreateStartLog(host string, port string) string {
	startMessage := fmt.Sprintf("LogCrunch server starting on %s:%s!", host, port)

	self, err := os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v", err)
		self = "this_machine"
	}

	startLog := structs.Log{
		Host:      self,
		Timestamp: time.Now().Unix(),
		Type:      "LogCrunch Server",
		Raw:       startMessage,
	}

	startLogJson, err := json.Marshal(startLog)
	if err != nil {
		log.Printf("Error marshalling startLog: %v", err)
		return "Failed to create startLog"
	}

	return string(startLogJson)
}
