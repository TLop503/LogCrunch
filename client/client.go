package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/TLop503/heartbeat0/cryptohandler"
	"github.com/TLop503/heartbeat0/structs"
	"github.com/joho/godotenv"
)

func makeHeartbeat(typ string, seq int) structs.Heartbeat {
	return structs.Heartbeat{
		Type:      typ,
		Timestamp: time.Now().Unix(),
		Seq:       seq,
	}
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

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

	encryptionKey := os.Getenv("AES_KEY")
	if len(encryptionKey) != 32 {
		fmt.Println("Invalid AES key length: must be 32 bytes")
		return
	}

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

		// Encrypt the JSON data
		encryptedData, err := cryptohandler.Encrypt(encryptionKey, string(jsonData))
		if err != nil {
			fmt.Println("Error encrypting heartbeat:", err)
			break
		}

		// Send the encrypted data with a newline
		_, err = writer.WriteString(encryptedData + "\n")
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			break
		}
		writer.Flush()
		// fmt.Printf("Sent (encrypted): %s\n", encryptedData)

		time.Sleep(5 * time.Second) // Send a heartbeat every 5 seconds
	}
}
