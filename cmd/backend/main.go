package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/piercegov/llm-npc-backend/internal/cfg"
)

func main() {
	config := cfg.ReadConfig()

	fmt.Printf("Starting server on port %s\n", config.Port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "LLM NPC Backend is running!")
	})

	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
