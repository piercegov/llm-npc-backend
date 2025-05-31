package main

import (
	"fmt"
	"net/http"
	"os"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logging.LogHTTPRequest(r.Method, r.URL.Path, http.StatusOK, "remote_addr", r.RemoteAddr)
		fmt.Fprintf(w, "LLM NPC Backend is running!")
	})
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logging.LogHTTPRequest(r.Method, r.URL.Path, http.StatusOK, "remote_addr", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})
	
	logging.Info("Server starting", "address", ":"+config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		logging.Error("Server failed to start", "error", err, "port", config.Port)
		os.Exit(1)
	}
}
