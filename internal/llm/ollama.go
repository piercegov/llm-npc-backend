package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/piercegov/llm-npc-backend/internal/cfg"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

// ollamaFunctionCall is a local struct to parse Ollama's specific tool call structure in responses.
type ollamaFunctionCall struct {
	Function struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	} `json:"function"`
}

// ollamaResponse represents the full JSON response from Ollama API
type ollamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role      string               `json:"role"`
		Content   string               `json:"content"`
		ToolCalls []ollamaFunctionCall `json:"tool_calls"`
	} `json:"message"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

// ollamaToolFunctionDetails defines the "function" part of a tool definition for requests.
type ollamaToolFunctionDetails struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ollamaTool defines the structure for a single tool in the request to Ollama.
type ollamaTool struct {
	Type     string                    `json:"type"` // Typically "function"
	Function ollamaToolFunctionDetails `json:"function"`
}

type Ollama struct {
	baseURL string
}

func NewOllama(baseURL string) *Ollama {
	return &Ollama{baseURL: baseURL}
}

func (o *Ollama) Generate(request LLMRequest) (LLMResponse, error) {
	// Transform Tool to Ollama-specific tool format
	var formattedTools []ollamaTool
	if len(request.Tools) > 0 {
		formattedTools = make([]ollamaTool, len(request.Tools))
		for i, t := range request.Tools {
			// Convert our Tool format to Ollama's expected format
			params := map[string]interface{}{
				"type":       "object",
				"properties": make(map[string]interface{}),
				"required":   []string{},
			}

			// Convert parameters
			for name, param := range t.Parameters {
				params["properties"].(map[string]interface{})[name] = map[string]interface{}{
					"type":        string(param.Type),
					"description": param.Description,
				}
				if param.Required {
					params["required"] = append(params["required"].([]string), name)
				}
			}

			formattedTools[i] = ollamaTool{
				Type: "function",
				Function: ollamaToolFunctionDetails{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  params,
				},
			}
		}
	}

	ollamaModel := cfg.ReadConfig().OllamaModel

	messages := []map[string]interface{}{}

	// Add system message if provided
	if request.SystemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": request.SystemPrompt,
		})
	}

	// Add user message
	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": request.Prompt,
	})

	requestMap := map[string]interface{}{
		"model":    ollamaModel,
		"messages": messages,
		"stream":   false,
	}

	// Only add the "tools" field if there are formatted tools
	if len(formattedTools) > 0 {
		requestMap["tools"] = formattedTools
	}

	jsonBody, err := json.Marshal(requestMap)
	if err != nil {
		logging.Error("Failed to marshal Ollama request body", "error", err, "requestMap_keys", func() []string {
			keys := make([]string, 0, len(requestMap))
			for k := range requestMap {
				keys = append(keys, k)
			}
			return keys
		}()) // Log keys of the map on error for diagnostics instead of the whole map
		return LLMResponse{}, err
	}

	// Log request in a more readable format
	logging.Info("Sending request to Ollama",
		"model", ollamaModel,
		"system_prompt_length", len(request.SystemPrompt),
		"user_prompt_length", len(request.Prompt),
		"tools_count", len(formattedTools),
	)
	logging.Debug("Ollama request details",
		"system_prompt", request.SystemPrompt,
		"user_prompt", request.Prompt,
	)

	httpRequest, err := http.NewRequest("POST", o.baseURL+"/api/chat", bytes.NewBuffer(jsonBody))
	if err != nil {
		logging.Error("Error creating request", "error", err)
		return LLMResponse{}, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	// Create client with configurable timeout
	config := cfg.ReadConfig()
	client := &http.Client{
		Timeout: config.LLMTimeout,
	}
	response, err := client.Do(httpRequest)

	if err != nil {
		logging.Error("Failed to send request to Ollama", "error", err)
		// Check if it's a timeout
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return LLMResponse{}, NewProviderError("ollama", ollamaModel, ErrTimeout, "request timed out")
		}
		return LLMResponse{}, NewProviderError("ollama", ollamaModel, ErrProviderUnavailable, "failed to connect to Ollama")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Error("Failed to read response body", "error", err)
		return LLMResponse{}, NewProviderError("ollama", ollamaModel, err, "failed to read response")
	}

	// Check HTTP status code
	if response.StatusCode != http.StatusOK {
		logging.Error("Ollama returned non-200 status",
			"status_code", response.StatusCode,
			"body", string(body),
		)

		// Map status codes to appropriate errors
		var baseErr error
		var message string
		switch response.StatusCode {
		case http.StatusBadRequest:
			baseErr = ErrBadRequest
			message = "invalid request parameters"
		case http.StatusUnauthorized:
			baseErr = ErrUnauthorized  
			message = "authentication failed"
		case http.StatusNotFound:
			baseErr = ErrModelNotFound
			message = fmt.Sprintf("model '%s' not found", ollamaModel)
		case http.StatusTooManyRequests:
			baseErr = ErrRateLimited
			message = "rate limit exceeded"
		case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
			baseErr = ErrProviderUnavailable
			message = "Ollama service unavailable"
		case http.StatusGatewayTimeout:
			baseErr = ErrTimeout
			message = "gateway timeout"
		default:
			baseErr = fmt.Errorf("unexpected status code: %d", response.StatusCode)
			message = string(body)
		}

		return LLMResponse{}, NewProviderError("ollama", ollamaModel, baseErr, message)
	}

	// Parse the full Ollama response
	var parsedResp ollamaResponse
	if err := json.Unmarshal(body, &parsedResp); err != nil {
		logging.Error("Failed to unmarshal Ollama response", "error", err, "body", string(body))
		return LLMResponse{}, NewProviderError("ollama", ollamaModel, err, "invalid response format")
	}

	// Extract tool uses from the parsed response
	toolUses := make([]ToolUse, len(parsedResp.Message.ToolCalls))
	for i, ollamaCall := range parsedResp.Message.ToolCalls {
		toolUses[i] = ToolUse{
			ToolName: ollamaCall.Function.Name,
			ToolArgs: ollamaCall.Function.Arguments,
		}
	}

	return LLMResponse{
		StatusCode: response.StatusCode,
		Response:   parsedResp.Message.Content, // Extract only the content
		ToolUses:   toolUses,
	}, nil
}
