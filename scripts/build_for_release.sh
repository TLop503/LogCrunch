#!/bin/bash

mkdir -p bins
mkdir -p release
go build -o bins/siem_intake_server ./server/siem_intake_server.go
go build -o bins/siem_agent ./agent/siem_agent.go

# Copy binaries and config file
cp bins/siem_intake_server release/siem_intake_server
cp bins/siem_agent release/siem_agent
cp targets.cfg release/targets.cfg

# Create a compressed tar archive
tar -czvf release.tar.gz release

# Cleanup release directory
rm -r release
