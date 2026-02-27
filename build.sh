#!/bin/bash

mkdir -p bins
<<<<<<< Updated upstream
mkdir -p release
go build -o bins/siem_intake_server ./server/siem_intake_server.go
go build -o bins/siem_agent ./agent/siem_agent.go

# Copy binaries and config file
cp bins/siem_intake_server release/siem_intake_server
cp bins/siem_agent release/siem_agent
cp targets.cfg release/targets.cfg
=======
go build -o bins/hb0-siem-intake-server ./server/server.go
go build -o bins/hb0-siem-agent ./agent/agent.go

# Create staging directory
mkdir -p staging/certs
mkdir -p staging/logs

# Generate new self-signed certificates
openssl req -x509 -newkey rsa:4096 -keyout staging/certs/server.key -out staging/certs/server.crt -days 365 -nodes -subj "/CN=localhost"

# Copy binaries and config file
cp bins/hb0-siem-intake-server staging/hb0-siem-intake-server
cp bins/hb0-siem-agent staging/hb0-siem-agent
cp targets.cfg staging/targets.cfg
>>>>>>> Stashed changes

# Create a compressed tar archive
tar -czvf release.tar.gz release

# Cleanup release directory
rm -r release
