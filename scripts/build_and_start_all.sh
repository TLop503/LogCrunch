#!/bin/bash

# For starting the server and client and manually triggering behaviors for testing
# intends to replace "scripts/start.sh"
# Should be called from root of the project

HOST="127.0.0.1"
PORT="5000"

# Build server and agent
mkdir -p bins
go build -o bins/siem_intake_server ./server
go build -o bins/siem_agent ./agent

# Create crypto IF it doesn't exist yet.
# Only for development
mkdir -p crypto
if [[ ! -f crypto/server.key || ! -f crypto/server.crt ]]; then
    echo "Generating development TLS certs..."
    openssl req -x509 -newkey rsa:4096 -keyout crypto/server.key -out crypto/server.crt -days 365 -nodes -subj "/CN=localhost"
else
    echo "TLS certs already exist, skipping generation."
fi

./bins/siem_intake_server $HOST $PORT crypto/server.crt crypto/server.key &
sleep 2
./bins/siem_agent $HOST $PORT targets.cfg n & # use default example targets file, and don't verify crypto
sleep 8
./bins/siem_agent $HOST $PORT targets.cfg n & # launch a second agent to observe behavior for duplicate hosts
