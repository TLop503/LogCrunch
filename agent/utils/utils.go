package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TLop503/LogCrunch/structs"
	"net"
	"os"
)

// GetHostName is a wrapper to handle the error-case of os.Hostname
func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return ("UnknownHost")
	}
	return hostname
}

// TransmitJson encodes json over a connection, reading inputs from a channel
func TransmitJson(conn net.Conn, logChan <-chan structs.Log) {
	for log := range logChan {
		err := json.NewEncoder(conn).Encode(log)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
	}
}
