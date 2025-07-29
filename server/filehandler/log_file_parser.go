package filehandler

import (
	"github.com/TLop503/LogCrunch/structs"
	"os"
)

// Path to where
const LOG_INTAKE_DESTINATION = "/var/log/LogCrunch.log"

func LogFileToData() (structs.IntakeLogFileData, error) {
	content, err := os.ReadFile(LOG_INTAKE_DESTINATION)
	if err != nil {
		return structs.IntakeLogFileData{}, err
	}

	return structs.IntakeLogFileData{FileContent: string(content)}, nil
}
