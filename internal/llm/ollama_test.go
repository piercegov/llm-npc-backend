package llm

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

	// Initialize Ollama with the mock server's URL
	ollama := NewOllama(server.URL)

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

	// Use the default Ollama URL for the integration test.
	ollama := NewOllama("http://localhost:11434")

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

// TestOllama_Generate_ModelNotFound tests the 404 model not found error
func TestOllama_Generate_ModelNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Return 404 Not Found
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"error": "model 'nonexistent-model' not found, try pulling it first"}`))
	}))
	defer server.Close()

	ollama := NewOllama(server.URL)
	_, err := ollama.Generate(LLMRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	// Check if the error is wrapped correctly
	var provErr *ProviderError
	if !errors.As(err, &provErr) {
		t.Fatalf("Expected ProviderError, got %T", err)
	}

	// Check if the underlying error is ErrModelNotFound
	if !errors.Is(err, ErrModelNotFound) {
		t.Errorf("Expected ErrModelNotFound, got %v", err)
	}
}

// TestOllama_Generate_BadRequest tests the 400 bad request error
func TestOllama_Generate_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"error": "invalid request format"}`))
	}))
	defer server.Close()

	ollama := NewOllama(server.URL)
	_, err := ollama.Generate(LLMRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !errors.Is(err, ErrBadRequest) {
		t.Errorf("Expected ErrBadRequest, got %v", err)
	}
}

// TestOllama_Generate_Unauthorized tests the 401 unauthorized error
func TestOllama_Generate_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(`{"error": "invalid API key"}`))
	}))
	defer server.Close()

	ollama := NewOllama(server.URL)
	_, err := ollama.Generate(LLMRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized, got %v", err)
	}
}

// TestOllama_Generate_RateLimited tests the 429 rate limit error
func TestOllama_Generate_RateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusTooManyRequests)
		rw.Write([]byte(`{"error": "rate limit exceeded"}`))
	}))
	defer server.Close()

	ollama := NewOllama(server.URL)
	_, err := ollama.Generate(LLMRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("Expected ErrRateLimited, got %v", err)
	}
}

// TestOllama_Generate_ServiceUnavailable tests the 503 service unavailable error
func TestOllama_Generate_ServiceUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusServiceUnavailable)
		rw.Write([]byte(`{"error": "service temporarily unavailable"}`))
	}))
	defer server.Close()

	ollama := NewOllama(server.URL)
	_, err := ollama.Generate(LLMRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !errors.Is(err, ErrProviderUnavailable) {
		t.Errorf("Expected ErrProviderUnavailable, got %v", err)
	}
}

// TestOllama_Generate_Timeout tests request timeout
func TestOllama_Generate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Simulate a slow response that will timeout
		time.Sleep(2 * time.Second)
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a custom Ollama instance with a very short timeout for testing
	// Note: In production, timeout is configured via cfg.ReadConfig()
	// ollama := NewOllama(server.URL)

	// We can't easily test the actual timeout without modifying the Generate method
	// This is a limitation of the current design
	// In a real test, you might want to add a way to inject the HTTP client
}

// TestOllama_Generate_ConnectionRefused tests connection refused error
func TestOllama_Generate_ConnectionRefused(t *testing.T) {
	// Use a port that's likely not in use
	ollama := NewOllama("http://localhost:54321")
	_, err := ollama.Generate(LLMRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Expected error but got none")
	}

	if !errors.Is(err, ErrProviderUnavailable) {
		t.Errorf("Expected ErrProviderUnavailable, got %v", err)
	}
}
