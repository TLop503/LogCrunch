package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/TLop503/LogCrunch/server/filehandler"
	"github.com/TLop503/LogCrunch/server/self_logging"
	web "github.com/TLop503/LogCrunch/server/web"
	"github.com/TLop503/LogCrunch/structs"
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

	log.Printf("TLS server listening on %s:%s\n", host, port)
	// log starting point
	filehandler.RotateFile("/var/log/LogCrunch/firehose.log",
		"/var/log/LogCrunch/old_firehose.log",
		true,
	)

	startLog := self_logging.CreateStartLog(host, port)
	filehandler.WriteToFile("/var/log/LogCrunch/firehose.log", true, false, startLog)

	// accept incoming transmissions indefinitely until we are killed
	connList := structs.NewConnList()
	// start web server
	web.Start(":8080", connList)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection:", err)
			continue
		}
		connList.AddToConnList(conn)
		go handleConnection(conn, connList)
	}
}

// takes an active connection and a pointer to the list of connections
// processes incoming logs (currently just writes to file)
// and updates the connection in the list when it is closed.
func handleConnection(conn net.Conn, connList *structs.ConnectionList) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())

	if err != nil {
		fmt.Println("Invalid remote address:", conn.RemoteAddr())
		return
	}

	hostNameSet := false
	hostname := ""

	for {
		agentFeedIn, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}

		if !hostNameSet {
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(agentFeedIn), &payload); err != nil {
				log.Println("Invalid JSON when scanning for hostname:", err)
			} else if val, ok := payload["host"]; ok {
				if hostStr, ok := val.(string); ok {
					hostname = hostStr
					hostNameSet = true
					fmt.Println("Hostname: ", hostname)
				}
			}
		}

		// Read connection from list
		// TODO: is mutex required here? could the connlist get away without one, since each conn has one?
		connList.RLock()
		trackedConn, ok := connList.Connections[host]
		connList.RUnlock()
		if ok {
			trackedConn.Lock()
			trackedConn.LastSeen = time.Now() // this should update after each received log entry.
			trackedConn.Hostname = hostname
			trackedConn.Unlock()
		}

		err = filehandler.WriteToFile(filehandler.LOG_INTAKE_DESTINATION, true, true, agentFeedIn)
		if err != nil {
			log.Println("Error writing file uncaught by file handler:", err)
		}
	}
}
