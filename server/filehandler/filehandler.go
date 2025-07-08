package filehandler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// WriteToFile writes a payload to a file, creating or appending based on flags.
// WriteToFile writes a payload to a file, creating or appending based on flags.
func WriteToFile(path string, create bool, append bool, payload string) error {
	var file *os.File
	var err error

	// Check if the file exists
	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		if create {
			// Create intermediate directories
			if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return fmt.Errorf("mkdir err: %v", err)
			}
			// Create the file if it doesn't exist and create flag is true
			file, err = os.Create(path)
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

		file, err = os.OpenFile(path, flags, 0644)

		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
	}

	defer file.Close()

	// Write the payload to the file
	_, err = file.WriteString(payload + "\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// RotateFile copies the contents of a file to a new location for archival purposes.
// filePath is the original file, rotationDestination is the new or existing file to write to
// append defines whether to append (true) or overwrite (false) any potentially existing data
// by default if the rotationDestination does not exist it will be created.
// Note! if file to rotate does not exist, this function just does nothing w/o error
func RotateFile(filePath string, rotationDestination string, append bool) error {
	// Check if the file exists
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		log.Printf("File (%s) to rotate does not exist!", filePath)
		return nil
	}

	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file (%s): %w", filePath, err)
	}

	err = WriteToFile(rotationDestination, true, append, string(contents))
	if err != nil {
		// since WriteToFile has verbose errors, we can just pass it upstream
		return fmt.Errorf("func WriteToFile error from RotateFile: %w", err)
	}

	return nil
}
