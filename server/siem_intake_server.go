package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	dbmod "github.com/TLop503/LogCrunch/server/db"
	"github.com/TLop503/LogCrunch/server/filehandler"
	"github.com/TLop503/LogCrunch/server/self_logging"
	"github.com/TLop503/LogCrunch/server/web"
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

	// initialize DB
	db, err := dbmod.InitDB("/var/log/LogCrunch/logcrunch.db")
	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}
	defer db.Close()
	roDB, err := sql.Open("sqlite3", "file:/var/log/LogCrunch/logcrunch.db?mode=ro&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("Error opening read-only DB: %v", err)
	}
	defer roDB.Close()
	err = dbmod.PrintAllModules(roDB)
	if err != nil {
		log.Fatalf("Error reading all modules in DB: %v", err)
	}

	// load modules from mpregistry
	err = dbmod.LoadModulesFromRegistry(db)
	if err != nil {
		log.Fatalf("Error loading modules: %v", err)
	}

	startLog := self_logging.CreateStartLog(host, port)
	err = filehandler.WriteToFile("/var/log/LogCrunch/firehose.log", true, false, startLog)
	if err != nil {
		log.Fatalf("Error initializing firehose: %v", err)
	}

	// accept incoming transmissions indefinitely until we are killed
	connList := structs.NewConnList()
	// start web server
	web.Start(":8080", connList, roDB) // use RO db connection!

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("Error accepting connection: %v", err)
			continue
		}
		connList.AddToConnList(conn)
		go handleConnection(conn, connList, db)
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

		err = dbmod.InsertLog(db, logStruct)
		if err != nil {
			log.Fatalf("Error inserting log into DB: %v", err)
		}
	}
}
