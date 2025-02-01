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

	//start the writer
	go writerRoutine(writer, logChan)

	// spin up a heartbeat goroutine to send proof of life
	// once every minute
	go heartbeat.Heartbeat(logChan, getHostName())
	go hemoglobin.ReadLog(logChan, "/var/log/auth.log")

	// TODO: Add graceful shutdowns
	select {}
}
