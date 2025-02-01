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

	// spin up a heartbeat goroutine to send proof of life
	// once every minute
	go heartbeat.Heartbeat(writer, getHostName())
	go hemoglobin.ReadLog("/var/log/auth.log", writer)

	// TODO: Add graceful shutdowns
	select {}
}
