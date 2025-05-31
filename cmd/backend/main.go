package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/piercegov/llm-npc-backend/internal/api"
	"github.com/piercegov/llm-npc-backend/internal/cfg"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

func main() {
	config := cfg.ReadConfig()
	
	// Initialize structured logging
	logging.InitLogger(config.LogLevel)
	
	logging.Info("Starting LLM NPC Backend server", 
		"port", config.Port, 
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

	// Apply middleware to handlers
	http.Handle("/", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(rootHandler, "GET"),
	))
	
	http.Handle("/health", api.ApplyDefaultMiddleware(
		api.WithMethodValidation(healthHandler, "GET"),
	))
	
	logging.Info("Server starting", "address", ":"+config.Port)
	
	// Set up the server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: nil, // Use default ServeMux
	}

	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		logging.Error("Server failed to start", "error", err, "port", config.Port)
		os.Exit(1)
	}
}
