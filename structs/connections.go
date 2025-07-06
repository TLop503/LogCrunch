package structs

import (
	"fmt"
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

func (ct *ConnectionList) AddToConnList(conn net.Conn) {
	ct.Lock()
	defer ct.Unlock()

	addr := conn.RemoteAddr().String()
	ct.Connections[addr] = &Connection{
		RemoteAddr: conn.RemoteAddr().String(),
		FirstSeen:  time.Now(),
		LastSeen:   time.Now(),
	}
	ct.print() // for testing, view all connections each time a new one arrives.
}

func (ct *ConnectionList) print() {
	fmt.Println("Active connections:")
	for _, conn := range ct.Connections {
		fmt.Printf("- %s / First seen: %s / Last Seen: %s\n", conn.RemoteAddr, conn.FirstSeen, conn.LastSeen)
	}
}
