package hemoglobin

import (
	"bufio"
	"fmt"
	"log"

	"github.com/hpcloud/tail"
)

/*
Hemoglobin is a <routine> containing <data> that facilitates the transportation of <logs> in <agents>.
*/

func ReadLog(path string, writer *bufio.Writer) {
	tailConfig := tail.Config{
		ReOpen:    true,                  // Handle log rotation
		Follow:    true,                  // Continuously read new lines
		MustExist: false,                 // Don't error if the file doesn't exist initially
		Logger:    tail.DiscardingLogger, // Disable internal logging
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
		_, err = writer.WriteString(line.Text)
		if err != nil {
			fmt.Println("Error sending log:", err)
			break
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing data:", err)
			break
		}
	}
}
