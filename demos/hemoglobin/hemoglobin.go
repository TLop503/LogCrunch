package main

import (
	"fmt"
	"log"

	"github.com/hpcloud/tail"
)

func authlog() {
	// Define the path to the log file
	logFile := "/var/log/auth.log" // Replace with the log file you want to monitor

	// Set tailing configuration
	tailConfig := tail.Config{
		ReOpen:    true,                  // Handle log rotation
		Follow:    true,                  // Continuously read new lines
		MustExist: false,                 // Don't error if the file doesn't exist initially
		Logger:    tail.DiscardingLogger, // Disable internal logging
	}

	// Open the log file with the specified configuration
	t, err := tail.TailFile(logFile, tailConfig)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	// Process lines as they are added to the log file
	for line := range t.Lines {
		if line.Err != nil {
			fmt.Printf("Error reading line: %v\n", line.Err)
			continue
		}
		// Print the log line
		fmt.Println(line.Text)
	}
}

func main() {
	go authlog()
}
