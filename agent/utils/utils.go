package utils

import (
	"bufio"
	"fmt"
	"os"
)

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return ("UnknownHost")
	}
	return hostname
}

// writerRoutine handles all writes to the server
func WriterRoutine(writer *bufio.Writer, dataChan <-chan string) {
	for data := range dataChan {
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			fmt.Println("Error writing data:", err)
			return
		}
		writer.Flush()
	}
}

// readTargets reads the target log file paths from ./targets.cfg
func ReadTargets(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			paths = append(paths, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return paths, nil
}
