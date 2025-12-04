#! /bin/bash

HOST="127.0.0.1"
PORT="5000"

# clean up from last time
./kill.sh

go run ./server/siem_intake_server.go $HOST $PORT ./certs/server.crt ./certs/server.key &
sleep 2
go run ./agent/siem_agent.go $HOST $PORT ./targets.cfg n &
