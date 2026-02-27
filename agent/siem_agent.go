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
<<<<<<< Updated upstream:agent/siem_agent.go
	if len(os.Args) < 5 {
		fmt.Println("Usage: program <host> <port> <congfig file> <verify certs y/n")
=======
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <host> <port> <InsecureSkipVerify (t/f)")
>>>>>>> Stashed changes:agent/agent.go
		return
	}

	fmt.Println("DISCLAIMER!!! This software ships with a default dummy cert bundled in for testing. DO NOT USE THIS CERT IN ANY REAL WORLD ENVIORNMENT; Supply YOUR OWN. This software is provided as-is and comes with no warranty!")
	fmt.Println(" __    _____  ___     ___  ____  __  __  _  _  ___  _   _ ")
	fmt.Println("(  )  (  _  )/ __)   / __)(  _ \\(  )(  )( \\( )/ __)( )_( )")
	fmt.Println(" )(__  )(_)(( (_-.  ( (__  )   / )(__)(  )  (( (__  ) _ ( ")
	fmt.Println("(____)(_____)\\___/   \\___)(_)\\_)(______)(_)\\_)\\___)(_) (_)")

	host := os.Args[1]
	port := os.Args[2]
<<<<<<< Updated upstream:agent/siem_agent.go
	cfg := os.Args[3]
	fmt.Println(os.Args[4])
	ISV := (os.Args[4] == "n")
	fmt.Println(ISV)
=======
	var ISV bool
	if os.Args[3] == "t" {
		ISV = true
	} else {
		ISV = false
	}
>>>>>>> Stashed changes:agent/agent.go

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
	targetPaths, err := utils.ReadTargets(cfg)
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
