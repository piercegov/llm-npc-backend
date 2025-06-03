# Unix Socket Server Testing Guide

This comprehensive guide covers how to build, run, and test the Unix socket server for the LLM NPC Backend.

## Prerequisites

- Go 1.22.0 or later installed
- Unix-like operating system (macOS, Linux)
- curl installed (for testing)

## Building the Server

1. **Clone and prepare the repository** (if not already done):
   ```bash
   git clone https://github.com/piercegov/llm-npc-backend.git
   cd llm-npc-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the server**:
   ```bash
   go build -o backend ./cmd/backend/...
   ```
   This creates an executable named `backend` in your current directory.

## Running the Server

### Using Default Socket Path

1. **Run with default configuration**:
   ```bash
   ./backend
   ```
   The server will use the default socket path: `/tmp/llm-npc-backend.sock`

2. **Verify the server started**:
   Look for log output similar to:
   ```
   {"time":"2025-01-06T10:30:00.000Z","level":"INFO","msg":"Starting LLM NPC Backend server","socket":"/tmp/llm-npc-backend.sock","log_level":"info","cerebras_base_url":"https://api.cerebras.ai"}
   {"time":"2025-01-06T10:30:00.001Z","level":"INFO","msg":"Server listening on Unix socket","socket":"/tmp/llm-npc-backend.sock"}
   ```

### Using Custom Socket Path

1. **Create a `.env` file** in the project root:
   ```bash
   echo "SOCKET_PATH=/tmp/my-custom-npc.sock" > .env
   ```

2. **Or use environment variable directly**:
   ```bash
   SOCKET_PATH=/tmp/my-custom-npc.sock ./backend
   ```

3. **Verify custom socket path in logs**:
   ```
   {"time":"2025-01-06T10:30:00.000Z","level":"INFO","msg":"Starting LLM NPC Backend server","socket":"/tmp/my-custom-npc.sock","log_level":"info","cerebras_base_url":"https://api.cerebras.ai"}
   ```

## Checking if the Socket Exists

1. **List the socket file**:
   ```bash
   ls -la /tmp/llm-npc-backend.sock
   ```
   You should see output like:
   ```
   srwxr-xr-x  1 user  wheel  0 Jan  6 10:30 /tmp/llm-npc-backend.sock
   ```
   Note: The `s` at the beginning indicates it's a socket file.

2. **Check if socket is in use**:
   ```bash
   lsof -U | grep llm-npc-backend.sock
   ```

## Testing with curl

### Basic Connectivity Test

1. **Test root endpoint**:
   ```bash
   curl --unix-socket /tmp/llm-npc-backend.sock http://localhost/
   ```
   Expected response:
   ```
   LLM NPC Backend is running!
   ```

2. **Test health endpoint**:
   ```bash
   curl --unix-socket /tmp/llm-npc-backend.sock http://localhost/health
   ```
   Expected response:
   ```
   pong
   ```

### Testing API Endpoints

1. **Test NPC endpoint** (requires Ollama running):
   ```bash
   curl --unix-socket /tmp/llm-npc-backend.sock http://localhost/npc
   ```
   This returns a JSON response with NPC data and LLM-generated content.

2. **Pretty-print JSON response**:
   ```bash
   curl -s --unix-socket /tmp/llm-npc-backend.sock http://localhost/npc | jq .
   ```

3. **Test with custom headers**:
   ```bash
   curl --unix-socket /tmp/llm-npc-backend.sock \
        -H "X-Request-ID: test-123" \
        http://localhost/health
   ```

### Testing Error Handling

1. **Test 404 Not Found**:
   ```bash
   curl --unix-socket /tmp/llm-npc-backend.sock http://localhost/nonexistent
   ```
   Expected: JSON error response with 404 status

2. **Test Method Not Allowed**:
   ```bash
   curl -X POST --unix-socket /tmp/llm-npc-backend.sock http://localhost/health
   ```
   Expected: JSON error response with 405 status

## Advanced Testing

### Using Different Log Levels

1. **Run with debug logging**:
   ```bash
   LOG_LEVEL=debug ./backend
   ```

2. **Test and observe detailed logs**:
   ```bash
   # In another terminal
   curl --unix-socket /tmp/llm-npc-backend.sock http://localhost/health
   ```
   You'll see detailed request/response logging.

### Testing with Custom Socket Permissions

1. **Create socket in a custom directory**:
   ```bash
   mkdir -p ~/npc-sockets
   SOCKET_PATH=~/npc-sockets/backend.sock ./backend
   ```

2. **Verify permissions**:
   ```bash
   ls -la ~/npc-sockets/backend.sock
   ```

### Load Testing

1. **Simple concurrent requests**:
   ```bash
   # Send 100 requests with 10 concurrent connections
   for i in {1..100}; do
     curl -s --unix-socket /tmp/llm-npc-backend.sock http://localhost/health &
   done
   wait
   ```

2. **Using Apache Bench (ab) with Unix sockets** (requires socat):
   ```bash
   # Create a TCP-to-Unix socket proxy
   socat TCP-LISTEN:8080,fork UNIX-CONNECT:/tmp/llm-npc-backend.sock &
   
   # Run load test
   ab -n 1000 -c 10 http://localhost:8080/health
   
   # Clean up
   killall socat
   ```

## Debugging Tips

### Common Issues and Solutions

1. **Socket file already exists**:
   ```bash
   # Error: bind: address already in use
   # Solution: Remove the old socket file
   rm /tmp/llm-npc-backend.sock
   ./backend
   ```

2. **Permission denied**:
   ```bash
   # Ensure you have write permissions to the socket directory
   ls -ld /tmp
   # Should show write permissions for your user
   ```

3. **Server not responding**:
   ```bash
   # Check if the process is running
   ps aux | grep backend
   
   # Check system logs
   tail -f /var/log/system.log  # macOS
   journalctl -u backend -f     # Linux with systemd
   ```

### Monitoring Socket Activity

1. **Watch socket file creation/deletion**:
   ```bash
   watch -n 1 'ls -la /tmp/*.sock 2>/dev/null || echo "No sockets found"'
   ```

2. **Monitor server logs in real-time**:
   ```bash
   ./backend 2>&1 | jq -R 'try fromjson catch .'
   ```
   This pretty-prints JSON log entries.

## Graceful Shutdown

1. **Send interrupt signal** (Ctrl+C):
   ```bash
   # The server will log:
   # {"level":"INFO","msg":"Shutting down server..."}
   ```

2. **Using kill command**:
   ```bash
   # Find the process ID
   ps aux | grep backend
   
   # Send SIGTERM
   kill -TERM <PID>
   ```

3. **Verify socket cleanup**:
   ```bash
   ls -la /tmp/llm-npc-backend.sock
   # Should return: No such file or directory
   ```

## Integration Testing

### Testing with Game Engine Clients

For game engines to connect to the Unix socket, they typically need Unix socket support in their HTTP client libraries. Here's a conceptual example:

```gdscript
# Godot example (conceptual - requires Unix socket support)
var client = HTTPClient.new()
client.connect_to_unix_socket("/tmp/llm-npc-backend.sock")
client.request(HTTPClient.METHOD_GET, "/health")
```

### Testing with Python Client

```python
import requests_unixsocket
import json

# Create a session that supports Unix sockets
session = requests_unixsocket.Session()

# Make requests using http+unix:// scheme
resp = session.get('http+unix://%2Ftmp%2Fllm-npc-backend.sock/health')
print(resp.text)  # Should print: pong

# Test NPC endpoint
resp = session.get('http+unix://%2Ftmp%2Fllm-npc-backend.sock/npc')
print(json.dumps(resp.json(), indent=2))
```

## Continuous Testing

### Shell Script for Health Monitoring

Create a file `monitor.sh`:
```bash
#!/bin/bash
SOCKET_PATH="${SOCKET_PATH:-/tmp/llm-npc-backend.sock}"

while true; do
    if [ -S "$SOCKET_PATH" ]; then
        RESPONSE=$(curl -s --unix-socket "$SOCKET_PATH" http://localhost/health)
        if [ "$RESPONSE" = "pong" ]; then
            echo "$(date): Server is healthy"
        else
            echo "$(date): Server returned unexpected response: $RESPONSE"
        fi
    else
        echo "$(date): Socket not found at $SOCKET_PATH"
    fi
    sleep 5
done
```

Make it executable and run:
```bash
chmod +x monitor.sh
./monitor.sh
```

## Summary

This testing guide covers:
- ✅ Building and running the server
- ✅ Testing with curl commands
- ✅ Checking if the socket exists
- ✅ Verifying the default socket path works
- ✅ Testing with a custom socket path
- ✅ Error handling and debugging
- ✅ Load testing and monitoring
- ✅ Integration testing examples

The Unix socket server provides efficient IPC communication for game engines, with easy testing through standard Unix tools like curl and comprehensive logging for debugging.