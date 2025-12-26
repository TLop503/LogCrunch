package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	userauth "github.com/TLop503/LogCrunch/server/user_auth"
	"log"
	"net"
	"os"
	"time"

	logdb "github.com/TLop503/LogCrunch/server/db/logs"
	"github.com/TLop503/LogCrunch/server/filehandler"
	"github.com/TLop503/LogCrunch/server/self_logging"
	"github.com/TLop503/LogCrunch/server/webserver"
	"github.com/TLop503/LogCrunch/structs"
)

func main() {

	fmt.Println("\n __    _____  ___     ___  ____  __  __  _  _  ___  _   _ ")
	fmt.Println("(  )  (  _  )/ __)   / __)(  _ \\(  )(  )( \\( )/ __)( )_( )")
	fmt.Println(" )(__  )(_)(( (_-.  ( (__  )   / )(__)(  )  (( (__  ) _ ( ")
	fmt.Println("(____)(_____)\\___/   \\___)(_)\\_)(______)(_)\\_)\\___)(_) (_)")

	if len(os.Args) < 5 {
		fmt.Println("Usage: <log_host> <log_port> <cert_path> <key_path> [http_host] [http_port]")
		fmt.Println("  log_host/log_port: Address for TLS log intake")
		fmt.Println("  cert_path/key_path: TLS certificate and key files")
		fmt.Println("  http_host/http_port: Address for webserver interface (optional, defaults to localhost:8080)")
		return
	}

	logHost := os.Args[1]
	logPort := os.Args[2]
	crt := os.Args[3]
	key := os.Args[4]

	// Default HTTP server settings
	httpHost := "localhost"
	httpPort := "8080"

	// Override HTTP settings if provided
	if len(os.Args) >= 6 {
		httpHost = os.Args[5]
	}
	if len(os.Args) >= 7 {
		httpPort = os.Args[6]
	}

	httpAddr := httpHost + ":" + httpPort

	// Load TLS certificate and key
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		log.Fatalf("Error loading TLS certificate and key: %v", err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", logHost+":"+logPort, config)
	if err != nil {
		log.Fatalf("Error starting TLS server: %v", err)
	}
	defer listener.Close()

	log.Printf("TLS server listening on %s:%s\n", logHost, logPort)
	// log starting point
	filehandler.RotateFile("/var/log/LogCrunch/firehose.log",
		"/var/log/LogCrunch/old_firehose.log",
		true,
	)

	// initialize log DBs
	logDB, roDB, err := logdb.InitLogDB("/var/log/LogCrunch/logcrunch.logDB")
	if err != nil {
		log.Fatalf("Error initializing DB connections: %v", err)
	}
	defer logDB.Close()
	defer roDB.Close()

	// initialize user database. create default user ad hoc
	userDB, err := userauth.FirstTimeSetupCheck("/opt/LogCrunch/users/accounts.userDB", "/opt/LogCrunch/users/.setupCompleted")
	defer userDB.Close()

	// TODO: pull out to 1-liner in self_logging
	startLog := self_logging.CreateStartLog(logHost, logPort)
	err = filehandler.WriteToFile("/var/log/LogCrunch/firehose.log", true, false, startLog)
	if err != nil {
		log.Fatalf("Error initializing firehose: %v", err)
	}

	// accept incoming transmissions indefinitely until we are killed
	connList := structs.NewConnList()
	// start webserver server
	webserver.StartRouter(httpAddr, connList, roDB) // use RO logDB connection!

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection: %v", err)
			continue
		}
		connList.AddToConnList(conn)
		go handleConnection(conn, connList, logDB)
	}
}

// takes an active connection and a pointer to the list of connections
// processes incoming logs (currently just writes to file)
// and updates the connection in the list when it is closed.
func handleConnection(conn net.Conn, connList *structs.ConnectionList, db *sql.DB) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		fmt.Println("Invalid remote address:", conn.RemoteAddr())
		return
	}

	hostNameSet := false
	hostname := ""

	for {
		var logEntry structs.Log
		if err := decoder.Decode(&logEntry); err != nil {
			if err.Error() == "EOF" {
				log.Println("Connection closed by remote")
			} else {
				log.Println("Failed to decode JSON:", err)
			}
			return
		}

		// Set hostname from the first log if not already set
		if !hostNameSet {
			hostname = logEntry.Host
			hostNameSet = true
		}

		// Update tracked connection info for webui
		connList.RLock()
		trackedConn, ok := connList.Connections[host]
		connList.RUnlock()
		if ok {
			trackedConn.Lock()
			trackedConn.LastSeen = time.Now()
			trackedConn.Hostname = hostname
			trackedConn.Unlock()
		}

		// Write raw JSON line to intake file
		// Currently kept in for debugging, may be deprecated in future.
		if err := filehandler.WriteToFile(filehandler.LOG_INTAKE_DESTINATION, true, true, logEntry); err != nil {
			log.Println("Error writing file uncaught by file handler:", err)
		}

		logStruct := structs.Log{
			Name:      logEntry.Name,
			Path:      logEntry.Path,
			Host:      logEntry.Host,
			Timestamp: logEntry.Timestamp,
			Module:    logEntry.Module,
			Parsed:    logEntry.Parsed,
			Raw:       logEntry.Raw,
		}

		err = logdb.InsertLog(db, logStruct)
		if err != nil {
			log.Fatalf("Error inserting log into DB: %v", err)
		}
	}
}
