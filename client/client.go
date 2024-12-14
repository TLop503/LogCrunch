package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Heartbeat struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Seq       int    `json:"seq"`
}

func makeHeartbeat(seq int) Heartbeat {
	return Heartbeat{
		Type:      "Heartbeat",
		Timestamp: time.Now().Unix(),
		Seq:       seq,
	}
}

func main() {
	host := "127.0.0.1"
	port := "5000"
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s:%s\n", host, port)
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	seq := 0

	for {
		// Create the heartbeat
		heartbeat := makeHeartbeat(seq)
		seq++

		// Convert the heartbeat to JSON
		jsonData, err := json.Marshal(heartbeat)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			break
		}

		// Send the JSON data with a newline
		_, err = writer.WriteString(string(jsonData) + "\n")
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			break
		}
		writer.Flush()
		fmt.Printf("Sent: %s\n", jsonData)

		// Read the response
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			break
		}
		fmt.Printf("Response: %s", response)

		time.Sleep(5 * time.Second) // Send a heartbeat every 5 seconds
	}
}
