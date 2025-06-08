# LLM NPC Backend API Documentation

This document describes the REST API endpoints for the LLM NPC Backend service.

## Base URL

- **Unix Socket Mode**: Communication via Unix domain socket at `/tmp/llm-npc-backend.sock`
- **HTTP Mode**: `http://localhost:8080` (default port, configurable via `HTTP_PORT` env var)

## Authentication

Currently, no authentication is required for API endpoints.

## Request Tracing

All requests are automatically traced with:
- Unique request ID (included in response headers and error responses)
- Response timing information
- Structured logging for debugging

## Response Format

### Success Response
Success responses vary by endpoint but typically include:
```json
{
  "success": true,
  "/* endpoint-specific data */": "..."
}
```

### Error Response
All endpoints use standardized error responses:
```json
{
  "error": "Human readable error message",
  "code": "ERROR_CODE",
  "request_id": "uuid-string (optional)",
  "details": {}
}
```

## Endpoints

### Health & Status

#### GET /
Root endpoint to verify the service is running.

**Response:**
```
LLM NPC Backend is running!
```

#### GET /health
Health check endpoint.

**Response:**
```
pong
```

### NPC Management

#### POST /npc/register
Register a new NPC with the system.

**Request Body:**
```json
{
  "name": "string (required)",
  "background_story": "string (required)"
}
```

**Response (201 Created):**
```json
{
  "npc_id": "uuid-string",
  "success": true,
  "message": "NPC registered successfully"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/npc/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Village Blacksmith",
    "background_story": "A skilled craftsman who has worked in this village for 20 years."
  }'
```

#### POST /npc/act
Execute a tick/action cycle for a specific NPC. This endpoint supports multi-round inference where NPCs can use tools and continue thinking for complex scenarios.

**Request Body:**
```json
{
  "npc_id": "string (required) - UUID of the NPC",
  "surroundings": [
    {
      "name": "string - Name of the surrounding element",
      "description": "string - Description of the element"
    }
  ],
  "knowledge_graph": {
    "nodes": [
      {
        "id": "string",
        "data": {}
      }
    ],
    "edges": [
      {
        "source": "string",
        "target": "string",
        "data": {}
      }
    ]
  },
  "npc_state": {},
  "knowledge_graph_depth": "integer (optional) - Depth for knowledge graph processing",
  "events": [
    {
      "event_type": "string - Type of event",
      "event_description": "string - Event description"
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "npc_id": "string",
  "rounds": [
    {
      "round_number": "integer",
      "llm_response": "string - NPC's response in this round",
      "tools_used": [
        {
          "tool_name": "string",
          "args": {},
          "success": "boolean",
          "response": "string",
          "error": "string"
        }
      ],
      "success": "boolean",
      "error_message": "string"
    }
  ],
  "llm_response": "string - Final NPC response",
  "success": "boolean",
  "error_message": "string"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/npc/act \
  -H "Content-Type: application/json" \
  -d '{
    "npc_id": "123e4567-e89b-12d3-a456-426614174000",
    "surroundings": [
      {
        "name": "forge",
        "description": "The forge is hot and busy. A customer approaches with a damaged sword."
      }
    ],
    "events": [
      {
        "event_type": "customer_arrival",
        "event_description": "A knight entered the shop carrying a broken sword"
      }
    ]
  }'
```

**Notes:**
- NPCs can execute up to 3 inference rounds if they use the `continue_thinking` tool
- Each round can include tool usage for complex decision-making
- The final `llm_response` contains the NPC's action or dialogue

#### GET /npc/list
List all registered NPCs.

**Response (200 OK):**
```json
{
  "npcs": {
    "uuid-string-1": {
      "name": "NPC Name",
      "background_story": "NPC background story"
    },
    "uuid-string-2": {
      "name": "Another NPC",
      "background_story": "Another background story"
    }
  },
  "success": true,
  "count": 2
}
```

#### GET /npc/{id}
Get detailed information for a specific NPC.

**Path Parameters:**
- `id` (string, required): UUID of the NPC

**Response (200 OK):**
```json
{
  "npc_id": "uuid-string",
  "npc": {
    "name": "NPC Name",
    "background_story": "NPC background story"
  },
  "success": true,
  "message": "NPC retrieved successfully"
}
```

#### DELETE /npc/{id}
Remove an NPC from the system.

**Path Parameters:**
- `id` (string, required): UUID of the NPC

**Response (200 OK):**
```json
{
  "npc_id": "uuid-string",
  "success": true,
  "message": "NPC deleted successfully"
}
```

### Development & Testing

#### GET /npc
Mock NPC endpoint for testing LLM integration.

**Note:** This is a temporary testing endpoint and may be removed in future versions.

**Response (200 OK):**
```json
{
  "npc_name": "Elara the Innkeeper",
  "background_story": "A warm and welcoming innkeeper who has been running 'The Gilded Swan' for over two decades...",
  "surroundings": [],
  "events": [],
  "llm_response": "Mock NPC response from LLM",
  "tools_used": [],
  "inference_rounds": 1,
  "tools_available": 2
}
```

#### GET /console/read_scratchpads
Development console endpoint for reading scratchpad data.

**Response (200 OK):**
```json
{
  "command": "read_scratchpads",
  "success": true,
  "data": {}
}
```

## Error Codes

| Code | Description |
|------|-------------|
| `INTERNAL_SERVER_ERROR` | Internal server error occurred |
| `INVALID_JSON` | Request body contains malformed JSON |
| `VALIDATION_ERROR` | Request validation failed (missing required fields) |
| `METHOD_NOT_ALLOWED` | HTTP method not allowed for this endpoint |
| `UNSUPPORTED_MEDIA_TYPE` | Content-Type header not supported |
| `RATE_LIMIT_EXCEEDED` | Too many requests |
| `SERVICE_UNAVAILABLE` | Service temporarily unavailable |
| `NOT_FOUND` | Resource not found |
| `BAD_REQUEST` | Invalid request parameters |

## Configuration

The API behavior can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `:8080` | HTTP server port (when using --http flag) |
| `SOCKET_PATH` | `/tmp/llm-npc-backend.sock` | Unix socket path |
| `OLLAMA_MODEL` | `qwen3:1.7b` | Ollama model to use |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |

## Usage Examples

### Complete NPC Workflow

1. **Register an NPC:**
```bash
curl -X POST http://localhost:8080/npc/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Innkeeper Sarah",
    "background_story": "A friendly innkeeper who knows all the local gossip and has run the Prancing Pony for 15 years."
  }'
```

2. **Execute NPC Actions:**
```bash
curl -X POST http://localhost:8080/npc/act \
  -H "Content-Type: application/json" \
  -d '{
    "npc_id": "your-npc-id-here",
    "surroundings": [
      {
        "name": "inn_main_room",
        "description": "The inn is quiet tonight. A hooded stranger sits in the corner nursing a drink."
      }
    ],
    "events": [
      {
        "event_type": "customer_behavior",
        "event_description": "The hooded stranger has been staring at the stairs leading to the rooms"
      }
    ]
  }'
```

3. **List NPCs:**
```bash
curl http://localhost:8080/npc/list
```

## Game Engine Integration

This API is designed to integrate with popular game engines:

- **Unity**: Use UnityWebRequest for HTTP calls
- **Godot**: Use HTTPRequest node
- **Unreal Engine**: Use HTTP module
- **Custom Engines**: Standard HTTP client libraries

### Recommended Integration Pattern

1. Register NPCs during game initialization using `/npc/register`
2. Call `/npc/act` during game tick cycles with current surroundings and events
3. Process NPC responses (including multi-round inference data) to update game state
4. Track events between ticks for temporal awareness
5. Use knowledge graphs for persistent NPC memory and world state

## Key Features

- **Multi-round Inference**: NPCs can think through complex scenarios over multiple rounds (max 3)
- **Tool Integration**: NPCs can execute tools for game-specific actions
- **Knowledge Graphs**: Support for persistent memory and world state representation
- **Request Tracing**: Every request includes unique ID and timing for debugging
- **Structured Logging**: Comprehensive logging with different levels

## Notes

- NPC IDs are UUIDs generated by the server
- The system supports both Unix socket and HTTP communication modes
- Enable debug logging (`LOG_LEVEL=debug`) to see full LLM prompts and responses
- NPCs use `<thinking>` tags for internal reasoning, separated from speech/actions
- Tool usage is tracked per inference round for detailed debugging