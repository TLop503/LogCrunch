#! /bin/bash

go run ./server/server.go &
sleep 2
go run ./agent/agent.go &
