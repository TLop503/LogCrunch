package structs

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Connection struct {
	sync.Mutex
	RemoteAddr string
	FirstSeen  time.Time
	LastSeen   time.Time
}

type ConnectionList struct {
	sync.RWMutex // RWMutex allows for better concurrent reads
	Connections  map[string]*Connection
}

// Create a new list of connections w/ initialized map
func NewConnList() *ConnectionList {
	return &ConnectionList{
		Connections: make(map[string]*Connection),
	}
}

// Add new host to list, w/ smart deduplication.
// duplicate connections are reconciled into a single entry in the connection list
func (ct *ConnectionList) AddToConnList(conn net.Conn) {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Println("Invalid remote address:", conn.RemoteAddr())
		return
	}

	ct.Lock()
	defer ct.Unlock()

	if existing, exists := ct.Connections[host]; exists {
		// Update metadata if already exists
		existing.Lock()
		existing.LastSeen = time.Now()
		existing.Unlock()
	} else {
		// Create a new entry
		ct.Connections[host] = &Connection{
			RemoteAddr: host,
			FirstSeen:  time.Now(),
			LastSeen:   time.Now(),
		}
	}

	ct.print()
}

func (ct *ConnectionList) print() {
	fmt.Println("Active connections:")
	for _, conn := range ct.Connections {
		fmt.Printf("- %s / First seen: %s / Last Seen: %s\n", conn.RemoteAddr, conn.FirstSeen, conn.LastSeen)
	}
}
