# AGENTS.md

This file guides Codex and other AI assistants working in this repository.

## Project Overview

This repository contains only the **Go backend** that acts as an LLM router for Non Player Characters in games. Each game engine (Unity, Godot, etc.) has its own lightweight connector script that communicates with this server over a Unix domain socket. The backend handles:

- Building prompts for NPCs
- Talking to different LLM providers
- Deciding which game actions an NPC can take (e.g. `throw`, `walk`)
- Returning structured results back to the client connector

## Repository Structure

```
cmd/backend/          - main.go entry point starting the socket server
internal/
  api/                - HTTP middleware and standardized errors
  cfg/                - configuration management via env vars
  kg/                 - knowledge graph structures
  llm/                - implementations of the LLMProvider interface
  logging/            - structured logging setup
  npc/                - core NPC logic and prompts
  store/              - (future) persistence layer
pkg/
  model/              - (placeholder for shared types)
```

See **CLAUDE.md** for a more exhaustive description of modules and design choices.
Keep this section in sync when new files are added.

## Development Commands

- Format Go code with `gofmt -w` before committing
- Run tests with `go test ./...`
- Build the server using `go build ./cmd/backend/...`
- During development you can run `go run ./cmd/backend`
- After building, start the executable directly with `./backend`

## Goal of the Backend

Act as a central router between game engine clients and Large Language Models. Game connectors should remain small, delegating prompt construction and NPC decision making entirely to this backend.
