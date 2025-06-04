# AGENTS.md

This file provides guidance to Codex and other AI assistants when working in this repository.

## Project Overview

This project is a **Go based backend** that serves as an LLM router for Non Player Characters in games. Game engines (for example Unity or Godot) will include small connector scripts that talk to this backend over a Unix domain socket. The backend is responsible for:

- Constructing prompts for the NPCs
- Interacting with different LLM providers
- Deciding which actions an NPC can take (e.g. `throw`, `walk`)
- Returning structured responses back to the game engine

Only the backend lives in this repository. Client/connector scripts will live in their own folders or repos and simply ping the server.

## Repository Structure

```
cmd/backend/          - main entry point that starts the socket server
internal/             - private packages
  api/                - HTTP middleware, error handling
  cfg/                - configuration management
  kg/                 - knowledge graph structures
  llm/                - implementations of the LLMProvider interface
  logging/            - structured logging setup
  npc/                - core NPC logic and prompts
  store/              - (future) persistence layer
pkg/                  - public packages
  model/              - (placeholder for shared types)
```

Keep this section up to date whenever new modules or files are added.

## Development Notes

- Format Go code using `gofmt -w` before committing.
- Run the full test suite with `go test ./...`.
- Build the server with `go build ./cmd/backend/...`.
- Quickly run the server during development with `go run ./cmd/backend`.
- After building, start the executable using `./backend` (or the generated path).

## Goal of the Backend

Act as a central router between game engine clients and Large Language Models. Engine connectors should be lightweight, delegating all LLM prompt building and NPC decision making to this backend.
