#! /bin/bash

HOST="127.0.0.1"
PORT="5000"

go run ./server/server.go $HOST $PORT &
sleep 2
go run ./agent/agent.go $HOST $PORT &
