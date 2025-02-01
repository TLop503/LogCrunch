package hemoglobin

import (
	"fmt"
	"log"

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
			fmt.Printf("Error reading line: %v\n", line.Err)
			continue
		}
		// write over wire
		//TODO! add parsing
		logChan <- line.Text
	}
}
