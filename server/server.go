package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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

func handleConnection(conn net.Conn) {
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

		// Process the heartbeat
		fmt.Printf("Processed Heartbeat: Type=%s, Timestamp=%d, Seq=%d\n", hb.Type, hb.Timestamp, hb.Seq)
		conn.Write([]byte("Heartbeat received and processed\n"))
	}
}
