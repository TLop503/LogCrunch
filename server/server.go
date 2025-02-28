package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/TLop503/heartbeat0/server/filehandler"
)

func main() {

	fmt.Println("\n __    _____  ___     ___  ____  __  __  _  _  ___  _   _ ")
	fmt.Println("(  )  (  _  )/ __)   / __)(  _ \\(  )(  )( \\( )/ __)( )_( )")
	fmt.Println(" )(__  )(_)(( (_-.  ( (__  )   / )(__)(  )  (( (__  ) _ ( ")
	fmt.Println("(____)(_____)\\___/   \\___)(_)\\_)(______)(_)\\_)\\___)(_) (_)")

	if len(os.Args) < 5 {
		fmt.Println("Usage: <host> <port> <cert path> <key path>")
		return
	}

	host := os.Args[1]
	port := os.Args[2]
	crt := os.Args[3]
	key := os.Args[4]

	// Load TLS certificate and key
	cert, err := tls.LoadX509KeyPair(crt, key)
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

	for {
		hb_in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}

		filehandler.WriteToFile("./logs/firehose.log", true, true, hb_in)
	}
}
