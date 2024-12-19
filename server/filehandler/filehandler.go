package filehandler

import (
	"fmt"
	"os"
)

// WriteToFile writes a payload to a file, creating or appending based on flags.
func WriteToFile(filePath string, create bool, append bool, payload string) error {
	var file *os.File
	var err error

	// Check if the file exists
	_, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		if create {
			// Create the file if it doesn't exist and create flag is true
			file, err = os.Create(filePath)
			if err != nil {
				return fmt.Errorf("error creating file: %w", err)
			}
		} else {
			// Error out if file doesn't exist and create flag is false
			return fmt.Errorf("file does not exist and create flag is false")
		}
	} else if err != nil {
		return fmt.Errorf("unexpected error checking file existence: %w", err)
	} else {
		//open file, and write or append
		flags := os.O_WRONLY
		if append {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		file, err = os.OpenFile(filePath, flags, 0644)

		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
	}

	defer file.Close()

	// Write the payload to the file
	_, err = file.WriteString(payload)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
