package hemoglobin

import (
	"encoding/json"
	"github.com/TLop503/LogCrunch/agent/utils"
	"github.com/TLop503/LogCrunch/structs"
	"log"
	"time"

	"github.com/hpcloud/tail"
)

/*
Hemoglobin is a <routine> containing <data> that facilitates the transportation of <logs> in <agents>.
*/

func ReadLog(logChan chan<- string, path string) {
	var seekOffset int64 = 0

	tailConfig := tail.Config{
		ReOpen:    true,                                          // Handle log rotation
		Follow:    true,                                          // Continuously read new lines
		MustExist: false,                                         // Don't error if the file doesn't exist initially
		Location:  &tail.SeekInfo{Offset: seekOffset, Whence: 0}, // Start from the end (TODO: Why doesn't this work)
		Logger:    tail.DiscardingLogger,                         // Disable internal logging
	}

	// Open the log file with the specified configuration
	t, err := tail.TailFile(path, tailConfig)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	// Process lines as they are added to the log file
	for line := range t.Lines {
		if line.Err != nil {
			log.Printf("Error reading line from file %v: %v\n", path, line.Err)
			continue
		}
		// write over wire
		//TODO! add parsing

		lcLog := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Type:      path,
			Raw:       line.Text,
		}
		jsonData, err := json.Marshal(lcLog)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
		}

		logChan <- string(jsonData)
	}
}
