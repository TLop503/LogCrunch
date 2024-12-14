package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Heartbeat struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Seq       int    `json:"seq"`
}

func main() {
	host := "127.0.0.1"
	port := "5000"
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s:%s\n", host, port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func fileHandler(createFlag bool, headMsg string, logJson string, tailMsg string) {
	var file *os.File
	var err error

	// Check if the file should be created
	if createFlag {
		file, err = os.Create("heartbeat.log") // Creates or truncates the file
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close() // Ensure the file is closed when done
	} else {
		// Open the file for appending if not creating
		file, err = os.OpenFile("heartbeat.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close() // Ensure the file is closed when done
	}

	// Write the bodyText to the file
	_, err = fmt.Fprintf(file, "%s: %s\n", headMsg, logJson, tailMsg)
	if err != nil {
		fmt.Println("Error writing to heartbeat.log:", err)
		return
	}
}

func handleConnection(conn net.Conn) {
	seq := 0
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read the incoming JSON message
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}
		fmt.Printf("Received: %s", message)

		// Parse the JSON into a struct
		var hb Heartbeat
		err = json.Unmarshal([]byte(message), &hb)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			conn.Write([]byte("Invalid JSON\n"))
			continue
		}

		// Check if the sequence number is out of order
		if hb.Seq != seq {
			// Marshal the heartbeat back into a JSON string to log it
			hbJson, err := json.Marshal(hb)
			if err != nil {
				fmt.Println("Error marshaling heartbeat to JSON:", err)
				continue
			}
			fileHandler(true, "SEQ out of order! Received heartbeat:", string(hbJson), fmt.Sprintf("Expected SEQ: %d", seq))
		} else {
			//update sequence
			seq++
		}

		// Process the heartbeat
		fmt.Printf("Processed Heartbeat: Type=%s, Timestamp=%d, Seq=%d\n", hb.Type, hb.Timestamp, hb.Seq)
		conn.Write([]byte("Heartbeat received and processed\n"))
	}
}
