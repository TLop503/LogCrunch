package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/TLop503/heartbeat0/structs"
)

func makeHeartbeat(typ string, seq int) structs.Heartbeat {
	ret := structs.Heartbeat{
		Type:      typ,
		Timestamp: time.Now().Unix(),
		Seq:       seq,
	}

	return ret
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

	seq := 9

	for {
		// Create the heartbeat
		heartbeat := makeHeartbeat("proof_of_life", seq)
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

		time.Sleep(5 * time.Second) // Send a heartbeat every 5 seconds
	}
}
