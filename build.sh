#!/bin/bash

HOST="127.0.0.1"
PORT="5000"

mkdir -p bins
go build -o bins/server ./server/server.go
go build -o bins/agent ./agent/agent.go

# Create staging directory
mkdir -p staging/certs
mkdir -p staging/logs

# Generate new self-signed certificates
openssl req -x509 -newkey rsa:4096 -keyout staging/certs/server.key -out staging/certs/server.crt -days 365 -nodes -subj "/CN=localhost"

# Copy binaries and config file
cp bins/server staging/server
cp bins/agent staging/agent
cp targets.cfg staging/targets.cfg

# Create a compressed tar archive
tar -czvf release.tar.gz staging

# Cleanup staging directory
rm -r staging
