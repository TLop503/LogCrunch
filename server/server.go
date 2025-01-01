package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/TLop503/heartbeat0/cryptohandler"
	"github.com/TLop503/heartbeat0/server/filehandler"
	"github.com/TLop503/heartbeat0/server/heartbeatlogs"
	"github.com/TLop503/heartbeat0/structs"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

	seq := 0
	encryptionKey := os.Getenv("AES_KEY")
	if len(encryptionKey) != 32 {
		log.Fatalf("Invalid AES key length: %d (must be 32 bytes for AES-256)", len(encryptionKey))
	}

	for {
		// Read the incoming JSON message
		hb_in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}

		// Decrypt the received log using cryptohandler
		plaintext, err := cryptohandler.Decrypt(encryptionKey, hb_in)
		if err != nil {
			log.Printf("Failed to decrypt log: %v", err)
			continue
		}

		// Decode decrypted JSON into the Heartbeat struct
		var hb structs.Heartbeat
		err = json.Unmarshal([]byte(plaintext), &hb)
		if err != nil {
			log.Printf("Failed to parse decrypted JSON: %v", err)
			continue
		}

		// Check for seq
		if hb.Seq != seq {
			// Seq error
			hblog, err := heartbeatlogs.GenerateSeqErrorLog("placeholder_host", seq, hb.Seq)
			if err != nil {
				log.Fatal(err)
			}
			filehandler.WriteToFile("heartbeat.log", true, true, hblog)
			seq = hb.Seq + 1 // After logging issue, reset seq
		} else {
			seq++
			hblog, err := heartbeatlogs.GenerateLog("placeholder_host", hb)
			if err != nil {
				log.Fatal(err)
			}
			filehandler.WriteToFile("heartbeat.log", true, true, hblog)
		}
	}
}
