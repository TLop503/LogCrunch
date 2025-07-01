#!/bin/bash
set -euo pipefail

# absolute path of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# relative paths
DOCKERFILE="$SCRIPT_DIR/dockerfiles/server_setup.Dockerfile"
BUILD_CONTEXT="$SCRIPT_DIR/../"
IMAGE_NAME="logcrunch-server-setup-test"
EXECUTABLE_PATH="/root/logcrunch_server" # inside the container!

echo "[*] Building Docker image..."
docker build -f "$DOCKERFILE" -t "$IMAGE_NAME" "$BUILD_CONTEXT"

echo "[*] Running container and simulating 'Enter' input..."
docker run --rm -i "$IMAGE_NAME" bash -c "
    echo '' | python3 /root/server_setup.py 5000 &&
    if [ -f \"$EXECUTABLE_PATH\" ]; then
        echo '[+] Success: logcrunch_server found'
        exit 0
    else
        echo '[!] Failure: logcrunch_server not found at $EXECUTABLE_PATH'
        exit 1
    fi
"
