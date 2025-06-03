#!/bin/bash
# Run server with debug logging to see network latency measurements

export LOG_LEVEL=debug
export SOCKET_PATH=/tmp/llm-npc-backend.sock

echo "Starting server with debug logging..."
echo "Socket path: $SOCKET_PATH"
echo "You'll see network latency measurements when clients send X-Client-Timestamp header"
echo ""

./backend