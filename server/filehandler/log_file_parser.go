package filehandler

import (
	"log"
	"os"

	"github.com/TLop503/LogCrunch/structs"
)

// Path to where
const LOG_INTAKE_DESTINATION = "/var/log/LogCrunch/firehose.log"

func LogFileToData() (*structs.IntakeLogFileData, error) {

	// Check if file exists first
	if _, err := os.Stat(LOG_INTAKE_DESTINATION); os.IsNotExist(err) {
		log.Printf("DEBUG: Log depot file does not exist: %s", LOG_INTAKE_DESTINATION)
		return &structs.IntakeLogFileData{FileContent: "Log file does not exist yet. No logs have been received."}, nil
	}

	content, err := os.ReadFile(LOG_INTAKE_DESTINATION)
	if err != nil {
		log.Printf("DEBUG: Error reading log depot file: %v", err)
		return &structs.IntakeLogFileData{}, err
	}

	log.Printf("DEBUG: Successfully read log depot file, content length: %d bytes", len(content))
	return &structs.IntakeLogFileData{FileContent: string(content)}, nil
}
