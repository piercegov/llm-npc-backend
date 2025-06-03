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
	"github.com/piercegov/llm-npc-backend/internal/llm"
	"github.com/piercegov/llm-npc-backend/internal/logging"
	"github.com/piercegov/llm-npc-backend/internal/npc"
)

func main() {
	// Initialize structured logging with default level
	logging.InitLogger("info")

	config := cfg.ReadConfig()

	// Reinitialize logger with configured log level
	logging.InitLogger(config.LogLevel)

	// Remove any existing socket file
	os.Remove(config.SocketPath)

	logging.Info("Starting LLM NPC Backend server",
		"socket", config.SocketPath,
		"log_level", config.LogLevel,
		"cerebras_base_url", config.BaseUrl)

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
			Events:              []npc.NPCTickEvent{}, // Empty events to see LLM behavior
		}

		// Parse the input to generate prompt
		surroundingsString, _ := npc.ParseSurroundings(mockInput)
		knowledgeGraphString, _ := npc.ParseKnowledgeGraph(mockInput)
		eventsString, _ := npc.ParseEvents(mockInput)

		systemPrompt := npc.BuildNPCSystemPrompt(mockNPC.Name, mockNPC.BackgroundStory)
		llmRequest := llm.LLMRequest{
			SystemPrompt: systemPrompt,
			Prompt:       surroundingsString + "\n" + knowledgeGraphString + "\n" + eventsString,
		}

		ollama := llm.NewOllama("11434")
		response, err := ollama.Generate(llmRequest)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate NPC response", api.ErrCodeInternalServer, nil, r.Context())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"npc_name":         mockNPC.Name,
			"background_story": mockNPC.BackgroundStory,
			"surroundings":     mockInput.Surroundings,
			"events":           mockInput.Events,
			"llm_response":     response.Response,
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
