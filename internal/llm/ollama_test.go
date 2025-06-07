package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllama_Generate_SuccessfulResponse(t *testing.T) {
	// Create a proper Ollama response structure
	expectedContent := "Hello! How can I help you today?"
	ollamaResponse := map[string]interface{}{
		"model":      "test-model",
		"created_at": "2023-01-01T00:00:00Z",
		"message": map[string]interface{}{
			"role":    "assistant",
			"content": expectedContent,
		},
		"done":        true,
		"done_reason": "stop",
	}
	jsonOllamaResponse, _ := json.Marshal(ollamaResponse)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Check if the request path is /api/chat
		if req.URL.Path != "/api/chat" {
			t.Errorf("Expected to request '/api/chat', got '%s'", req.URL.Path)
		}
		// Check if the request method is POST
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST request, got '%s'", req.Method)
		}

		// Send a 200 OK response with a proper Ollama JSON body
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonOllamaResponse)
	}))
	defer server.Close()

	// Extract port from server.URL (e.g., "http://127.0.0.1:12345" -> "12345")
	// The server.URL is like "http://127.0.0.1:PORT"
	urlParts := strings.Split(server.URL, ":")
	port := urlParts[len(urlParts)-1]

	// Initialize Ollama with the mock server's port
	ollama := NewOllama(port)

	prompt := "Hello, Ollama!"
	response, err := ollama.Generate(LLMRequest{Prompt: prompt})

	if err != nil {
		t.Fatalf("Generate() returned an unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if response.Response != expectedContent {
		t.Errorf("Generate() response content = %s, want %s", response.Response, expectedContent)
	}
}

// TestOllama_Generate_WithToolCall verifies that the response contains at least one tool call.
func TestOllama_Generate_WithToolCall(t *testing.T) {
	// Skip this test if -short is passed, as it's an integration test.
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	// Use the default Ollama port for the integration test.
	ollama := NewOllama("11434")

	// Prompt designed to trigger a tool call
	prompt := "Please use a tool to get the current weather in Paris in celsius."
	llmResponse, err := ollama.Generate(LLMRequest{Prompt: prompt, Tools: []Tool{makeWeatherTool()}})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			t.Skipf("Skipping integration test: Ollama instance not reachable at port 11434 (connection refused). Error: %v", err)
			return
		}
		t.Fatalf("Generate() returned an unexpected error: %v", err)
	}

	if llmResponse.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusOK, llmResponse.StatusCode, llmResponse.Response)
	}

	// Verify that the response contains at least one tool call
	if len(llmResponse.ToolUses) == 0 {
		t.Errorf("Expected at least one tool call in the response, but got none. Response: %s", llmResponse.Response)
	}
	//Validate that it called the get_current_weather tool with the correct arguments
	if llmResponse.ToolUses[0].ToolName != "get_current_weather" {
		t.Errorf("Expected tool name 'get_current_weather', got '%s'", llmResponse.ToolUses[0].ToolName)
	}
	if llmResponse.ToolUses[0].ToolArgs["location"] != "Paris" {
		t.Errorf("Expected location 'Paris', got '%s'", llmResponse.ToolUses[0].ToolArgs["location"])
	}
	if llmResponse.ToolUses[0].ToolArgs["format"] != "celsius" {
		t.Errorf("Expected format 'celsius', got '%s'", llmResponse.ToolUses[0].ToolArgs["format"])
	}

	t.Logf("Integration test received response with tool call: %s", llmResponse.Response)
}

func makeWeatherTool() Tool {
	return Tool{
		Name:        "get_current_weather",
		Description: "Get the current weather for a location",
		Parameters: map[string]ToolParameter{
			"location": {
				Type:        TypeString,
				Description: "The location to get the weather for, e.g. San Francisco, CA",
				Required:    true,
			},
			"format": {
				Type:        TypeString,
				Description: "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
				Required:    true,
			},
		},
	}
}
