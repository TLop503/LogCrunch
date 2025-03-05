package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/TLop503/LogCrunch/agent/heartbeat"
	"github.com/TLop503/LogCrunch/agent/hemoglobin"
	"github.com/TLop503/LogCrunch/agent/utils"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: program <host> <port> <verify certs y/n")
		return
	}

	host := os.Args[1]
	port := os.Args[2]
	ISV := os.Args[3] == "y"

	// Configure TLS
	config := &tls.Config{InsecureSkipVerify: ISV} // Set to `false` in production with valid certs
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
	go utils.WriterRoutine(writer, logChan)

	// spin up a heartbeat goroutine to send proof of life
	// once every minute
	go heartbeat.Heartbeat(logChan, utils.GetHostName())

	// Read log file paths from targets.cfg
	targetPaths, err := utils.ReadTargets("./targets.cfg")
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
