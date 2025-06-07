package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/piercegov/llm-npc-backend/internal/api"
	"github.com/piercegov/llm-npc-backend/internal/cfg"
	"github.com/piercegov/llm-npc-backend/internal/kg"
	"github.com/piercegov/llm-npc-backend/internal/logging"
	"github.com/piercegov/llm-npc-backend/internal/npc"
	"github.com/piercegov/llm-npc-backend/internal/tools"
)

// getAllToolsUsed extracts all tools used across all inference rounds
func getAllToolsUsed(rounds []npc.InferenceRound) []npc.ToolResult {
	var allTools []npc.ToolResult
	for _, round := range rounds {
		allTools = append(allTools, round.ToolsUsed...)
	}
	return allTools
}

func main() {
	// Initialize structured logging with default level
	logging.InitLogger("info")

	config := cfg.ReadConfig()

	// Reinitialize logger with configured log level
	logging.InitLogger(config.LogLevel)

	// Remove any existing socket file
	os.Remove(config.SocketPath)

	// Initialize tool registry and scratchpad storage
	toolRegistry := tools.NewToolRegistry()
	scratchpadStorage := tools.NewScratchpadStorage()

	// Register scratchpad tools
	if err := tools.RegisterScratchpadTools(toolRegistry, scratchpadStorage); err != nil {
		logging.Error("Failed to register scratchpad tools", "error", err)
		os.Exit(1)
	}

	// Initialize NPC storage and handlers
	npcStorage := npc.NewNPCStorage()
	npcHandlers := npc.NewNPCHandlers(npcStorage, toolRegistry)

	logging.Info("Starting LLM NPC Backend server",
		"socket", config.SocketPath,
		"log_level", config.LogLevel,
		"cerebras_base_url", config.BaseUrl,
		"tools_count", len(toolRegistry.GetTools()),
		"npc_endpoints", "POST /npc/register, POST /npc/act, GET /npc/list, GET /npc/{id}, DELETE /npc/{id}")

	// Define the root handler
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			api.WriteErrorResponse(w, http.StatusNotFound, "Not found", api.ErrCodeNotFound, nil, r.Context())
			return
		}
		fmt.Fprintf(w, "LLM NPC Backend is running!")
	})

	// Define the health check handler
	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})

	// Define the NPC handler
	npcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockNPC := npc.NPC{
			Name:            "Elara the Innkeeper",
			BackgroundStory: "A friendly innkeeper who has run the Prancing Pony tavern for over 20 years. She knows everyone in town and loves to share local gossip.",
		}

		mockInput := npc.NPCTickInput{
			Surroundings: []npc.Surrounding{
				{
					Name:        "Tavern Common Room",
					Description: "A warm, dimly lit room with wooden tables and chairs. Several patrons are drinking ale and chatting quietly.",
				},
				{
					Name:        "Stranger",
					Description: "A hooded figure sits alone in the corner, nursing a drink and watching the room carefully.",
				},
			},
			KnowledgeGraph:      kg.KnowledgeGraph{},
			NPCState:            npc.NPCState{},
			KnowledgeGraphDepth: 0,
			Events: []npc.NPCTickEvent{
				{
					EventType:        "new_customer",
					EventDescription: "The hooded stranger entered the tavern just moments ago and ordered a whiskey",
				},
			},
			ToolRegistry: toolRegistry, // Add the tool registry
		}

		// Use ActForTick which now returns detailed results
		result := mockNPC.ActForTick(mockInput)

		if !result.Success {
			api.WriteErrorResponse(w, http.StatusInternalServerError, result.ErrorMessage, api.ErrCodeInternalServer, nil, r.Context())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"npc_name":         mockNPC.Name,
			"background_story": mockNPC.BackgroundStory,
			"surroundings":     mockInput.Surroundings,
			"events":           mockInput.Events,
			"llm_response":     result.LLMResponse,
			"tools_used":       getAllToolsUsed(result.Rounds),
			"inference_rounds": len(result.Rounds),
			"tools_available":  len(toolRegistry.GetTools()),
		})
	})

	// Define the console handler for reading scratchpads
	consoleHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get all scratchpads from storage
		allScratchpads := scratchpadStorage.GetAllScratchpads()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"command": "read_scratchpads",
			"success": true,
			"data":    allScratchpads,
		})
	})

	// Apply middleware to handlers
	http.Handle("/", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(rootHandler, "GET"),
	))

	http.Handle("/health", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(healthHandler, "GET"),
	))

	http.Handle("/npc", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(npcHandler, "GET"),
	))

	http.Handle("/console/read_scratchpads", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(consoleHandler, "GET"),
	))

	// NPC management endpoints
	http.Handle("/npc/register", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(http.HandlerFunc(npcHandlers.RegisterHandler), "POST"),
	))

	http.Handle("/npc/act", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(http.HandlerFunc(npcHandlers.ActHandler), "POST"),
	))

	http.Handle("/npc/list", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(http.HandlerFunc(npcHandlers.ListHandler), "GET"),
	))

	// NPC-specific endpoints (GET and DELETE /npc/{id})
	npcGetDeleteHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			npcHandlers.GetHandler(w, r)
		case "DELETE":
			npcHandlers.DeleteHandler(w, r)
		default:
			api.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", api.ErrCodeMethodNotAllowed, nil, r.Context())
		}
	})

	http.Handle("/npc/", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(npcGetDeleteHandler, "GET", "DELETE"),
	))

	// Create Unix socket listener
	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		logging.Error("Failed to create Unix socket", "error", err, "socket", config.SocketPath)
		os.Exit(1)
	}
	defer listener.Close()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logging.Info("Shutting down server...")
		listener.Close()
		os.Remove(config.SocketPath)
	}()

	logging.Info("Server listening on Unix socket", "socket", config.SocketPath)

	// Start serving on the Unix socket
	err = http.Serve(listener, nil)
	if err != nil && err != net.ErrClosed {
		logging.Error("Server error", "error", err)
		os.Exit(1)
	}
}
