#! /bin/bash

go run ./server/server.go &
sleep 2
go run ./client/client.go &
