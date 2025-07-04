# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based backend for powering Non-Player Characters (NPCs) in video games using Large Language Models. The framework supports both Unix domain sockets for efficient IPC communication and HTTP mode for easier testing and integration with popular game engines (Godot, Unity, Unreal Engine).

## Development Commands

### Building and Running
- **Build**: `go build ./cmd/backend/...`
- **Run (Unix Socket)**: `./backend` (default mode, after building)
- **Run (HTTP)**: `./backend --http` (HTTP mode for easier testing/integration)
- **Dependencies**: `go mod tidy`

### Log Viewer
- **GUI Mode (Unix Socket)**: `./run_logviewer.sh` (default)
- **GUI Mode (HTTP)**: `./run_logviewer.sh --http`
- **CLI Mode (Unix Socket)**: `./run_logviewer.sh --cli`
- **CLI Mode (HTTP)**: `./run_logviewer.sh --cli --http`

### Testing
- **Run all tests**: `go test ./...`
- **Run specific package tests**: `go test ./internal/api`
- **Run tests with verbose output**: `go test -v ./...`
- **Unix socket testing**: See [TESTING_GUIDE.md](./TESTING_GUIDE.md) for comprehensive testing instructions

## Architecture

### Core Modules
- **`cmd/backend/main.go`**: Unix socket server entry point with basic health check and root endpoints
- **`internal/npc/`**: Core NPC logic with tick-based action system (`ActForTick`)
- **`internal/llm/`**: LLM provider interface and implementations (Ollama, LM Studio)
- **`internal/api/`**: HTTP middleware, error handling, and request tracing
- **`internal/kg/`**: Knowledge graph integration for NPC context and memory
- **`internal/cfg/`**: Configuration management via environment variables

### Key Concepts
- **NPCs**: Defined by `Name` and `BackgroundStory`, process `Surroundings` on each tick
- **NPCTickEvent**: Tracks events that occurred since the last tick for temporal awareness
- **LLM Provider Interface**: Abstraction layer supporting multiple LLM providers
- **Knowledge Graph**: Configurable depth system for NPC context and decision-making
- **Tool System**: LLMs can use predefined tools for game-specific actions
- **Tick-based Actions**: NPCs operate on `ActForTick` cycles for dynamic behavior

### Testing Patterns
- Uses standard Go testing with `httptest` for HTTP handlers
- Middleware tests validate request tracing, panic recovery, and method validation
- Tests use table-driven test patterns for multiple scenarios

## Configuration

Environment variables (can be set in `.env` file):
- `SOCKET_PATH`: Unix socket path (default: /tmp/llm-npc-backend.sock)
- `HTTP_PORT`: HTTP server port when using --http flag (default: :8080)
- `LOG_LEVEL`: Logging level - debug, info, warn, error (default: info)
  - Set to `debug` to see full LLM prompts and responses in logs

### LLM Provider Configuration
- `LLM_PROVIDER`: Selects the LLM provider - ollama, lmstudio (default: ollama)

#### Ollama Settings
- `OLLAMA_MODEL`: Ollama model to use (default: qwen3:1.7b)
- `OLLAMA_BASE_URL`: Ollama server URL (default: http://10.0.0.85:11434)

#### LM Studio Settings
- `LMSTUDIO_BASE_URL`: LM Studio server URL (default: http://localhost:1234)
- `LMSTUDIO_MODEL`: Model identifier for LM Studio (default: model)
- `LMSTUDIO_API_KEY`: API key for LM Studio (default: lm-studio)

#### Other Providers (Planned)
- `CEREBRAS_API_KEY`: Cerebras API key (optional)
- `CEREBRAS_BASE_URL`: Cerebras API base URL (default: https://api.cerebras.ai)

## Current API Endpoints

### Core Endpoints
- `GET /`: Root endpoint returning "LLM NPC Backend is running!"
- `GET /health`: Health check returning "pong"

### NPC Management Endpoints (Game Engine Integration)
- `POST /npc/register`: Register a new NPC with name and background story, returns NPC ID
- `POST /npc/act`: Execute a tick for a specific NPC using NPCTickInput data
- `GET /npc/list`: List all registered NPCs with their basic information
- `GET /npc/{id}`: Get detailed information for a specific NPC
- `DELETE /npc/{id}`: Remove an NPC from the system

### Development/Testing Endpoints
- `GET /npc`: Mock NPC endpoint demonstrating LLM integration (temporary testing endpoint)
- `GET /console/read_scratchpads`: Development console for reading scratchpad data

## Project Structure Guide

**IMPORTANT**: Always update this section when adding new files or changing the project structure.

### Directory Layout
```
llm-npc-backend/
├── cmd/
│   └── backend/
│       └── main.go              # Unix socket server entry point, route definitions, handlers
├── internal/                    # Private application code (not importable by other projects)
│   ├── api/
│   │   ├── errors.go           # Standardized error responses and error codes
│   │   ├── middleware.go       # HTTP middleware (logging, panic recovery, tracing, CORS)
│   │   └── middleware_test.go  # Middleware unit tests
│   ├── cfg/
│   │   └── cfg.go              # Configuration management, env var parsing
│   ├── kg/
│   │   └── kg.go               # Knowledge graph data structures (Node, Edge, KnowledgeGraph)
│   ├── llm/
│   │   ├── common.go           # LLM interfaces (LLMProvider, LLMRequest, LLMResponse, Tool)
│   │   ├── ollama.go           # Ollama LLM provider implementation
│   │   └── ollama_test.go      # Ollama provider tests
│   ├── logging/
│   │   └── logger.go           # Structured logging setup using slog
│   ├── npc/
│   │   ├── npc.go              # Core NPC logic, ActForTick, prompt parsing, API request/response types
│   │   ├── npc_test.go         # NPC unit tests
│   │   ├── prompts.go          # NPC system prompts and templates
│   │   ├── storage.go          # In-memory NPC storage with UUID-based identification
│   │   └── handlers.go         # HTTP handlers for NPC management endpoints
│   └── store/                  # (Placeholder for future data persistence)
├── pkg/                        # Public packages (can be imported by other projects)
│   └── model/                  # (Placeholder for shared data models)
├── go.mod                      # Go module definition
├── README.md                   # Project documentation for users
└── CLAUDE.md                   # This file - AI assistant guidance

### Module Responsibilities

#### `cmd/backend/main.go`
- Server initialization and configuration (Unix socket or HTTP mode)
- Command line flag parsing (--http flag for HTTP mode)
- Route definitions and handler setup
- Middleware application
- Currently contains mock NPC handler for testing

#### `internal/api/`
- **errors.go**: Defines error codes and standardized JSON error responses
- **middleware.go**: Implements cross-cutting concerns (request ID, logging, panic recovery, CORS)

#### `internal/cfg/`
- **cfg.go**: Central configuration management, reads environment variables, provides defaults

#### `internal/kg/`
- **kg.go**: Knowledge graph structures for NPC memory and world state representation

#### `internal/llm/`
- **common.go**: Defines provider-agnostic interfaces for LLM integration
  - `LLMProvider`: Interface for LLM implementations
  - `LLMRequest`: Includes SystemPrompt, Prompt, and Tools
  - `LLMResponse`: Contains response text and tool uses
  - `Tool`: Defines callable tools/functions for LLMs
- **ollama.go**: Ollama-specific implementation of LLMProvider
  - Properly parses Ollama API responses to extract only content
  - Structured logging instead of raw JSON for better readability
- **lmstudio.go**: LM Studio implementation of LLMProvider
  - Supports OpenAI-compatible API endpoints
  - Can connect to LM Studio instances on the local network
  - Handles tool calling and structured responses
- **factory.go**: Provider factory for creating LLM instances based on configuration
  - Dynamically selects between Ollama and LM Studio providers
  - Centralizes provider instantiation logic

#### `internal/logging/`
- **logger.go**: Configures structured logging with different log levels

#### `internal/npc/`
- **npc.go**: Core NPC behavior logic and API data structures
  - `ActForTick`: Main NPC action loop
  - Surroundings and knowledge graph parsing
  - Event tracking with `NPCTickEvent` for temporal awareness
  - LLM interaction with structured XML prompts
  - `ParseEvents`: Formats events since last tick in XML
  - API request/response types for game engine integration
- **storage.go**: Thread-safe in-memory NPC storage system
  - UUID-based NPC identification
  - CRUD operations with mutex protection
  - Registration, retrieval, listing, and deletion of NPCs
- **handlers.go**: HTTP handlers for NPC management endpoints
  - `POST /npc/register`: Register new NPCs
  - `POST /npc/act`: Execute NPC ticks with full input processing
  - `GET /npc/list`, `GET /npc/{id}`, `DELETE /npc/{id}`: Management operations
- **prompts.go**: System prompt templates
  - Instructs LLMs to use `<thinking>` tags for internal reasoning
  - Ensures clean separation of thoughts vs. speech/actions

### Key Design Patterns

1. **Interface-based Design**: LLM providers implement a common interface for easy swapping
2. **Middleware Chain**: Request/response concerns handled through composable middleware
3. **Structured Logging**: Consistent log format with contextual information
4. **Environment-based Config**: All configuration through environment variables
5. **Tick-based NPCs**: Game loop compatible design for NPC actions

### Adding New Features

When adding new features:
1. Place internal code in `internal/` with appropriate module
2. Update this structure guide with new files/modules
3. Add tests alongside implementation files
4. Update configuration in `cfg.go` if new env vars needed
5. Document new API endpoints in this file

### Testing Strategy
- Unit tests live alongside implementation files (*_test.go)
- Use table-driven tests for multiple scenarios
- Mock external dependencies (LLMs, databases)
- Integration tests should use the mock handlers