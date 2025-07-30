package hemoglobin

import (
	"github.com/TLop503/LogCrunch/agent/utils"
	"github.com/TLop503/LogCrunch/structs"
	"log"
	"time"

	"github.com/hpcloud/tail"
)

/*
Hemoglobin is a <routine> containing <data> that facilitates the transportation of <logs> in <agents>.
*/

func ReadLog(logChan chan<- structs.Log, path string) {
	var seekOffset int64 = 0

	tailConfig := tail.Config{
		ReOpen:    true,                                          // Handle log rotation
		Follow:    true,                                          // Continuously read new lines
		MustExist: false,                                         // Don't error if the file doesn't exist initially
		Location:  &tail.SeekInfo{Offset: seekOffset, Whence: 0}, // Start from the given offset (0 = beginning)
		Logger:    tail.DiscardingLogger,                         // Disable internal logging
	}

	t, err := tail.TailFile(path, tailConfig)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	for line := range t.Lines {
		if line.Err != nil {
			log.Printf("Error reading line from file %v: %v\n", path, line.Err)
			continue
		}

		logEntry := structs.Log{
			Host:      utils.GetHostName(),
			Timestamp: time.Now().Unix(),
			Type:      "TODO",
			Path:      path,
			Raw:       line.Text,
			Parsed:    nil,
		}

		// send to channel for writing across wire
		logChan <- logEntry
	}
}
