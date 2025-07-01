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
    echo '[*] Starting server setup...'
    start_time=\$(date +%s)
    
    echo '' | python3 /root/server_setup.py localhost 5000 &&
    
    end_time=\$(date +%s)
    duration=\$((end_time - start_time))
    echo \"[*] Setup script completed in \${duration} seconds\"
    
    # Check if binary exists
    if [ -f \"$EXECUTABLE_PATH\" ]; then
        echo '[+] Success: logcrunch_server binary found'
    else
        echo '[!] Failure: logcrunch_server not found at $EXECUTABLE_PATH'
        exit 1
    fi
    
    # Wait a moment for the server to start
    sleep 3
    
    # Check if server process is running
    if pgrep -f 'logcrunch_server' > /dev/null; then
        echo '[+] Success: logcrunch_server process is running'
        
        # Get process details
        echo '[*] Server process info:'
        ps aux | grep logcrunch_server | grep -v grep
        
        # Check if server is listening on port 5000
        echo '[*] Checking if server is listening on port 5000...'
        if ss -tlnp | grep ':5000 '; then
            echo '[+] Success: Server is listening on port 5000 (ss)'
    
        else
            echo '[!] Warning: Server process running but not detected listening on port 5000'
            echo '[*] All listening ports:'
            ss -tlnp 2>/dev/null
        fi
        
        # Terminate the server process for clean exit
        pkill -f 'logcrunch_server' || true
        sleep 1
        
        exit 0
    else
        echo '[!] Failure: logcrunch_server process is not running'
        echo '[*] Checking for any processes:'
        ps aux | head -10
        exit 1
    fi
"
