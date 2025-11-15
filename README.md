# LLM NPC Framework Backend

[![Documentation](https://img.shields.io/badge/docs-live-brightgreen?style=for-the-badge&logo=read-the-docs)](https://piercegov.github.io/llm-npc-backend/)
[![Python SDK](https://img.shields.io/badge/Python_SDK-available-blue?style=for-the-badge&logo=python)](https://piercegov.github.io/llm-npc-backend/)

This project provides a backend framework for powering Non-Player Characters (NPCs) in video games using Large Language Models (LLMs). It is designed to be compatible with popular game engines such as Godot, Unity, and Unreal Engine.

## Features
*   **LLM Agnostic Design:** Utilizes a flexible `LLMProvider` interface, enabling support for various Large Language Models.
    *   Currently implemented for Ollama and LM Studio.
    *   Planned support for other providers like OpenAI, Cerebras, Groq, etc.
*   **Unix Socket Communication:** Uses Unix domain sockets for efficient, low-latency IPC communication instead of HTTP.
*   **Core NPC Functionality:**
    *   NPCs are defined with a `Name` and `BackgroundStory`.
    *   Supports stateful NPCs with a (planned) configurable `NPCState` (e.g., for health, inventory, faction alignment).
    *   NPCs perceive their `Surroundings` (objects/characters with names and descriptions).
    *   Operates on a tick-based action cycle (`ActForTick`) for dynamic behavior.
    *   Tracks events since last tick via `NPCTickEvent` for temporal awareness.
*   **Knowledge Graph (KG) Integration:** NPCs can leverage an internal knowledge graph for richer context, memory, and decision-making. The depth of KG information used can be configured per interaction.
*   **Tool Usage by LLMs:** Supports prompting LLMs to use predefined "tools." This allows NPCs to trigger game-specific actions, perform lookups, or execute other programmed capabilities. Tools are defined with a name, description, and parameters.
*   **Configuration via Environment:** Easily configurable through environment variables or a `.env` file (socket path, LLM model selection, API keys, logging level).

## Supported Game Engines
*   Godot
*   Unity
*   Unreal Engine

## Architecture Overview
The backend is a Go-based server that communicates via Unix domain sockets. It receives requests from game engines, processes them through NPC logic modules, interacts with a configured Large Language Model (via the `LLMProvider` interface), and can utilize a knowledge graph for enhanced NPC responses and actions. Key interactions include NPC perception of surroundings, state updates, and LLM-driven actions which can involve tool usage.

(Further details or a diagram could be added here if desired)

## Getting Started

### Prerequisites
*   Go (version 1.22.0 or later recommended)
*   One of the following LLM providers:
    *   Ollama installed and running (if using the default Ollama provider)
    *   LM Studio installed and running with "Serve on Local Network" enabled (if using LM Studio)
*   Git

### Installation & Setup
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/piercegov/llm-npc-backend.git
    cd llm-npc-backend
    ```
2.  **Ensure Dependencies are Synced:**
    Run `go mod tidy` to synchronize your project's dependencies. This will download any necessary modules and clean up the `go.mod` and `go.sum` files.
    ```bash
    go mod tidy
    ```
3.  **Set up Environment Variables:**
    Create a `.env` file in the root of the project or set the environment variables directly. See the "Environment Variables" section for details.
    
    A minimal `.env` for Ollama might look like:
    ```env
    SOCKET_PATH=/tmp/llm-npc-backend.sock
    LLM_PROVIDER=ollama
    OLLAMA_MODEL=qwen3:1.7b # Or your preferred Ollama model
    LOG_LEVEL=info
    ```
    
    For LM Studio on a remote machine:
    ```env
    SOCKET_PATH=/tmp/llm-npc-backend.sock
    LLM_PROVIDER=lmstudio
    LMSTUDIO_BASE_URL=http://192.168.1.100:1234  # Replace with your LM Studio server IP
    LMSTUDIO_MODEL=llama-3.2-1b-instruct         # Replace with your loaded model
    LOG_LEVEL=info
    ```
4.  **Build the backend:**
    ```bash
    go build ./cmd/backend/...
    ```
5.  **Run the backend:**
    ```bash
    ./backend # Or the path to your built executable
    ```
    The server will start listening on the Unix socket (default: `/tmp/llm-npc-backend.sock`).

## Environment Variables
The following environment variables are used for configuration. You can set them directly or place them in a `.env` file in the project root.

### General Configuration
*   `SOCKET_PATH`: The Unix socket path where the server will listen.
    *   Default: `/tmp/llm-npc-backend.sock`
*   `LOG_LEVEL`: Sets the logging level for the application.
    *   Options: `debug`, `info`, `warn`, `error`
    *   Default: `info`
    *   Note: Use `debug` to see full LLM prompts and responses

### LLM Provider Selection
*   `LLM_PROVIDER`: Selects which LLM provider to use.
    *   Options: `ollama`, `lmstudio` (or `lm-studio`)
    *   Default: `ollama`

### Ollama Configuration
*   `OLLAMA_MODEL`: Specifies the model to be used with the Ollama provider.
    *   Default: `qwen3:1.7b`
    *   Example: `llama3:latest`, `mistral:latest`

### LM Studio Configuration
*   `LMSTUDIO_BASE_URL`: The base URL for the LM Studio server.
    *   Default: `http://localhost:1234`
    *   Example for remote machine: `http://192.168.1.100:1234`
*   `LMSTUDIO_MODEL`: The model identifier to use with LM Studio.
    *   Default: `model`
    *   Example: `llama-3.2-1b-instruct`, `deepseek-r1-distill-qwen-7b`
*   `LMSTUDIO_API_KEY`: API key for LM Studio (if required).
    *   Default: `lm-studio`

### Other Providers (Planned)
*   `CEREBRAS_API_KEY`: Your API key for the Cerebras LLM provider (if you plan to use it).
    *   Optional.
*   `CEREBRAS_BASE_URL`: The base URL for the Cerebras API.
    *   Default: `https://api.cerebras.ai`
    *   Optional.

## Usage / API Endpoints
The primary interaction with the backend is via Unix socket communication using HTTP protocol.

*   **`GET /`**: Root endpoint. Returns a simple message indicating the backend is running.
*   **`GET /health`**: Health check endpoint. Returns `pong` with a 200 OK status if the server is healthy.
*   **`GET /npc`**: Mock NPC endpoint for testing LLM integration (temporary).

For comprehensive testing instructions including curl commands, socket verification, and debugging tips, see the [Testing Guide](./TESTING_GUIDE.md).
