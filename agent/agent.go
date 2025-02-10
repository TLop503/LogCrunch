package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/TLop503/heartbeat0/agent/heartbeat"
	"github.com/TLop503/heartbeat0/agent/hemoglobin"
)

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return ("UnknownHost")
	}
	return hostname
}

// writerRoutine handles all writes to the server
func writerRoutine(writer *bufio.Writer, dataChan <-chan string) {
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
func readTargets(filePath string) ([]string, error) {
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

func main() {
	host := "127.0.0.1"
	port := "5000"

	// Configure TLS
	config := &tls.Config{InsecureSkipVerify: true} // Set to `false` in production with valid certs
	// Connect to server
	conn, err := tls.Dial("tcp", host+":"+port, config)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	fmt.Printf("Connected to %s:%s via TLS\n", host, port)

	// create channel for thread-safe writes
	logChan := make(chan string)

	// start the writer
	go writerRoutine(writer, logChan)

	// spin up a heartbeat goroutine to send proof of life
	// once every minute
	go heartbeat.Heartbeat(logChan, getHostName())

	// Read log file paths from targets.cfg
	targetPaths, err := readTargets("./targets.cfg")
	if err != nil {
		fmt.Println("Error reading targets file:", err)
		return
	}

	// Start a hemoglobin instance for each target path
	for _, path := range targetPaths {
		go hemoglobin.ReadLog(logChan, path)
	}

	// TODO: Add graceful shutdowns
	select {}
}
