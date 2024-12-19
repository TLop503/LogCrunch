package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/TLop503/heartbeat0/server/filehandler"
	"github.com/TLop503/heartbeat0/server/heartbeaterror"
	"github.com/TLop503/heartbeat0/structs"
)

func main() {
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
	var last_ts int64
	last_ts = 0

	for {
		// Read the incoming JSON message
		hb_in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}

		//decode to json
		var hb structs.Heartbeat
		err = json.Unmarshal([]byte(hb_in), &hb)
		if err != nil {
			log.Fatal(err)
		}

		//if first time, set last ts to expected
		last_ts = hb.Timestamp

		//check for seq
		if hb.Seq != seq {
			//seq error
			hblog, err := heartbeaterror.GenerateSeqErrorLog("placeholder_host", seq, hb.Seq)
			if err != nil {
				log.Fatal(err)
			}
			filehandler.WriteToFile("heartbeat.log", true, true, hblog)
			seq = hb.Seq + 1 //after logging issue reset seq
		} else {
			fmt.Printf("%+v\n", hb)
			seq++
		}

	}
}
