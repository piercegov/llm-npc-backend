package llm

import (
	"bytes"
	"encoding/json"
	"io"
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

// ollamaResponseForToolParsing is a helper struct to unmarshal the JSON response
// and access the nested tool_calls from Ollama.
type ollamaResponseForToolParsing struct {
	Message struct {
		ToolCalls []ollamaFunctionCall `json:"tool_calls"`
	} `json:"message"`
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
	OLLAMA_PORT string
}

func NewOllama(ollamaPort string) *Ollama {
	return &Ollama{OLLAMA_PORT: ollamaPort}
}

// extractToolUses parses the Ollama response and converts its tool calls
// to the common llm.ToolUse type.
func extractToolUses(responseBody string) []ToolUse { // Return type is the common llm.ToolUse
	var parsedOllamaResp ollamaResponseForToolParsing
	if err := json.Unmarshal([]byte(responseBody), &parsedOllamaResp); err != nil {
		logging.Error("Failed to unmarshal Ollama response for tool extraction", "error", err, "body", responseBody)
		return []ToolUse{} // Return empty slice of common llm.ToolUse on error
	}

	if parsedOllamaResp.Message.ToolCalls == nil {
		return []ToolUse{}
	}

	// Convert []ollamaFunctionCall to []llm.ToolUse
	commonToolUses := make([]ToolUse, len(parsedOllamaResp.Message.ToolCalls))
	for i, ollamaCall := range parsedOllamaResp.Message.ToolCalls {
		commonToolUses[i] = ToolUse{ // This now refers to llm.ToolUse from common.go
			ToolName: ollamaCall.Function.Name,
			ToolArgs: ollamaCall.Function.Arguments,
		}
	}
	return commonToolUses
}

func (o *Ollama) Generate(request LLMRequest) (LLMResponse, error) {
	// Transform common.Tool to Ollama-specific tool format
	var formattedTools []ollamaTool
	if len(request.Tools) > 0 { // Check if there are any tools to format
		formattedTools = make([]ollamaTool, len(request.Tools))
		for i, t := range request.Tools { // t is of type common.Tool (or the llm.Tool type from this package)
			formattedTools[i] = ollamaTool{
				Type:     "function", // Assuming all tools are functions for Ollama
				Function: ollamaToolFunctionDetails(t),
			}
		}
	}

	ollamaModel := cfg.ReadConfig().OllamaModel

	requestMap := map[string]interface{}{
		"model": ollamaModel,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": request.Prompt,
			},
		},
		"stream": false,
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

	// Restore original logging line
	logging.Info("Sending request to Ollama", "request", string(jsonBody))

	httpRequest, err := http.NewRequest("POST", "http://localhost:"+o.OLLAMA_PORT+"/api/chat", bytes.NewBuffer(jsonBody))
	if err != nil {
		logging.Error("Error creating request", "error", err)
		return LLMResponse{}, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(httpRequest)

	if err != nil {
		return LLMResponse{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return LLMResponse{}, err
	}

	toolUses := extractToolUses(string(body))

	return LLMResponse{
		StatusCode: response.StatusCode,
		Response:   string(body),
		ToolUses:   toolUses,
	}, nil
}
