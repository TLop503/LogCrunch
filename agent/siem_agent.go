package main

import (
	"crypto/tls"
	"fmt"
	"github.com/TLop503/LogCrunch/structs"
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/TLop503/LogCrunch/agent/heartbeat"
	"github.com/TLop503/LogCrunch/agent/hemoglobin"
	"github.com/TLop503/LogCrunch/agent/utils"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: program <host> <port> <congfig file> <verify certs y/n")
		return
	}

	host := os.Args[1]
	port := os.Args[2]
	cfg := os.Args[3]
	fmt.Println(os.Args[4])
	ISV := (os.Args[4] == "n")
	fmt.Println(ISV)

	// Configure TLS
	config := &tls.Config{InsecureSkipVerify: ISV} // Set to `false` in production with valid certs
	// Connect to server
	conn, err := tls.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to %s:%s via TLS\n", host, port)

	// create channel for thread-safe writes
	logChan := make(chan structs.Log)

	// start the writer
	go utils.TransmitJson(conn, logChan)

	// spin up a heartbeat goroutine to send proof of life
	// once every minute
	go heartbeat.Heartbeat(logChan, utils.GetHostName())

	// Read log file paths from config file`
	data, err := os.ReadFile(cfg)
	if err != nil {
		fmt.Errorf("Error reading targets file:", err)
		return
	}

	var yamlConfig structs.YamlConfig

	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		fmt.Errorf("Error unmarshalling targets file:", err)
		return
	}

	// Start a hemoglobin instance for each target path
	for _, target := range yamlConfig.Targets {
		go hemoglobin.ReadLog(logChan, target)
	}

	// TODO: Add graceful shutdowns
	select {}
}
