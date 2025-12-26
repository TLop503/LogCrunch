package filehandler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// WriteToFile writes a payload to a file, creating or appending based on flags.
// payload can be a string or any value that can be marshaled to JSON.
func WriteToFile(path string, create bool, append bool, payload interface{}) error {
	var file *os.File
	var err error

	// Marshal payload to string
	var line string
	switch v := payload.(type) {
	case string:
		line = v
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal payload to JSON: %w", err)
		}
		line = string(data)
	}

	// Check if the file exists
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if create {
			if err = Create_if_needed(path, os.ModePerm, os.ModePerm); err != nil {
				return fmt.Errorf("file/path creation error: %v", err)
			}
		} else {
			return fmt.Errorf("file does not exist and create flag is false")
		}
	} else if err != nil {
		return fmt.Errorf("unexpected error checking file existence: %w", err)
	} else {
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

	// Write the payload (JSON string) to the file
	_, err = file.WriteString(line + "\n")
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

// Create_if_needed creates an empty file at a given path to aid in DB or filehandler operations
func Create_if_needed(path string, dirPerm os.FileMode, filePerm os.FileMode) error {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Ensure DB file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, filePerm)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		f.Close()
	} else if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	return nil
}
