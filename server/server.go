package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/TLop503/heartbeat0/server/filehandler"
)

func main() {

	host := "127.0.0.1"
	port := "5000"

	// Load TLS certificate and key
	cert, err := tls.LoadX509KeyPair("./server/server.crt", "./server/server.key")
	if err != nil {
		log.Fatalf("Error loading TLS certificate and key: %v", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", host+":"+port, config)
	if err != nil {
		log.Fatalf("Error starting TLS server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("TLS server listening on %s:%s\n", host, port)

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

	//seq := 0

	for {
		// Read the incoming JSON message
		hb_in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}

		/* for testing
		// Decode JSON into the Heartbeat struct
		var hb structs.Heartbeat
		err = json.Unmarshal([]byte(hb_in), &hb)
		if err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			continue
		}

		// Check for seq
		if hb.Seq != seq {
			// Seq error
			hblog, err := heartbeatlogs.GenerateSeqErrorLog("placeholder_host", seq, hb.Seq)
			if err != nil {
				log.Fatal(err)
			}
			filehandler.WriteToFile("./logs/heartbeat.log", true, true, hblog)
			seq = hb.Seq + 1 // After logging issue, reset seq
		} else {
			seq++
			hblog, err := heartbeatlogs.GenerateLog(hb)
			if err != nil {
				log.Fatal(err)
			}
			filehandler.WriteToFile("./logs/heartbeat.log", true, true, hblog)
		}
		*/
		filehandler.WriteToFile("./logs/firehose.log", true, true, hb_in)
	}
}
