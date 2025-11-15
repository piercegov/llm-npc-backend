package llm

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLMStudioGenerate(t *testing.T) {
	// Create a test server to mock LM Studio API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("Expected path /v1/chat/completions, got %s", r.URL.Path)
		}

		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			t.Errorf("Expected Authorization header with Bearer token, got %s", authHeader)
		}

		// Parse request body
		var req lmStudioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Verify request contents
		if req.Model != "test-model" {
			t.Errorf("Expected model 'test-model', got %s", req.Model)
		}

		if len(req.Messages) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(req.Messages))
		}

		// Send mock response
		response := lmStudioResponse{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "test-model",
			Choices: []lmStudioChoice{
				{
					Index: 0,
					Message: struct {
						Role      string             `json:"role"`
						Content   string             `json:"content"`
						ToolCalls []lmStudioToolCall `json:"tool_calls,omitempty"`
					}{
						Role:    "assistant",
						Content: "Test response from LM Studio",
					},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create LM Studio provider with test server URL
	provider := NewLMStudio(server.URL, "test-model", "test-api-key")

	// Create test request
	request := LLMRequest{
		SystemPrompt: "You are a helpful assistant.",
		Prompt:       "Hello, world!",
	}

	// Generate response
	response, err := provider.Generate(request)
	if err != nil {
		t.Fatalf("Failed to generate response: %v", err)
	}

	// Verify response
	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	if response.Response != "Test response from LM Studio" {
		t.Errorf("Expected response 'Test response from LM Studio', got %s", response.Response)
	}
}

func TestLMStudioGenerateWithTools(t *testing.T) {
	// Create a test server to mock LM Studio API with tool calls
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request to verify tools
		var req lmStudioRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}

		// Verify tools were included
		if len(req.Tools) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(req.Tools))
		}

		// Send response with tool call
		response := lmStudioResponse{
			ID:      "chatcmpl-test-tools",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "test-model",
			Choices: []lmStudioChoice{
				{
					Index: 0,
					Message: struct {
						Role      string             `json:"role"`
						Content   string             `json:"content"`
						ToolCalls []lmStudioToolCall `json:"tool_calls,omitempty"`
					}{
						Role:    "assistant",
						Content: "",
						ToolCalls: []lmStudioToolCall{
							{
								ID:   "call_123",
								Type: "function",
								Function: struct {
									Name      string `json:"name"`
									Arguments string `json:"arguments"`
								}{
									Name:      "test_tool",
									Arguments: `{"param1": "value1"}`,
								},
							},
						},
					},
					FinishReason: "tool_calls",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create provider
	provider := NewLMStudio(server.URL, "test-model", "test-api-key")

	// Create request with tools
	request := LLMRequest{
		SystemPrompt: "You are a helpful assistant.",
		Prompt:       "Use the test tool.",
		Tools: []Tool{
			{
				Name:        "test_tool",
				Description: "A test tool",
				Parameters: map[string]ToolParameter{
					"param1": {
						Type:        TypeString,
						Description: "A test parameter",
						Required:    true,
					},
				},
			},
		},
	}

	// Generate response
	response, err := provider.Generate(request)
	if err != nil {
		t.Fatalf("Failed to generate response: %v", err)
	}

	// Verify tool uses
	if len(response.ToolUses) != 1 {
		t.Errorf("Expected 1 tool use, got %d", len(response.ToolUses))
	}

	if response.ToolUses[0].ToolName != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got %s", response.ToolUses[0].ToolName)
	}

	if response.ToolUses[0].ToolArgs["param1"] != "value1" {
		t.Errorf("Expected param1='value1', got %v", response.ToolUses[0].ToolArgs["param1"])
	}
}

func TestLMStudioGenerateError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	// Create provider
	provider := NewLMStudio(server.URL, "test-model", "test-api-key")

	// Create request
	request := LLMRequest{
		Prompt: "Test",
	}

	// Generate response
	_, err := provider.Generate(request)

	// Should get an error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that it's the right error type
	if !errors.Is(err, ErrProviderUnavailable) {
		t.Errorf("Expected ErrProviderUnavailable, got %v", err)
	}
}
