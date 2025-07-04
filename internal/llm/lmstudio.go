package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/piercegov/llm-npc-backend/internal/logging"
)

// lmStudioMessage represents a message in the LM Studio chat format
type lmStudioMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// lmStudioToolFunction represents the function part of a tool for LM Studio
type lmStudioToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// lmStudioTool represents a tool definition for LM Studio
type lmStudioTool struct {
	Type     string               `json:"type"`
	Function lmStudioToolFunction `json:"function"`
}

// lmStudioRequest represents the request payload for LM Studio API
type lmStudioRequest struct {
	Model       string            `json:"model"`
	Messages    []lmStudioMessage `json:"messages"`
	Tools       []lmStudioTool    `json:"tools,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	Stream      bool              `json:"stream"`
}

// lmStudioToolCall represents a tool call in the response
type lmStudioToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// lmStudioChoice represents a choice in the response
type lmStudioChoice struct {
	Index   int `json:"index"`
	Message struct {
		Role      string             `json:"role"`
		Content   string             `json:"content"`
		ToolCalls []lmStudioToolCall `json:"tool_calls,omitempty"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

// lmStudioResponse represents the response from LM Studio API
type lmStudioResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []lmStudioChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// LMStudio implements the LLMProvider interface for LM Studio
type LMStudio struct {
	BaseURL string
	Model   string
	APIKey  string
}

// NewLMStudio creates a new LM Studio provider instance
func NewLMStudio(baseURL, model, apiKey string) *LMStudio {
	if apiKey == "" {
		apiKey = "lm-studio" // Default API key for LM Studio
	}
	return &LMStudio{
		BaseURL: baseURL,
		Model:   model,
		APIKey:  apiKey,
	}
}

// Generate implements the LLMProvider interface
func (l *LMStudio) Generate(request LLMRequest) (LLMResponse, error) {
	// Build messages array
	messages := []lmStudioMessage{}

	// Add system message if provided
	if request.SystemPrompt != "" {
		messages = append(messages, lmStudioMessage{
			Role:    "system",
			Content: request.SystemPrompt,
		})
	}

	// Add user message
	messages = append(messages, lmStudioMessage{
		Role:    "user",
		Content: request.Prompt,
	})

	// Convert tools to LM Studio format
	var tools []lmStudioTool
	if len(request.Tools) > 0 {
		tools = make([]lmStudioTool, len(request.Tools))
		for i, tool := range request.Tools {
			// Build parameters schema
			params := map[string]interface{}{
				"type":       "object",
				"properties": make(map[string]interface{}),
				"required":   []string{},
			}

			// Add parameter definitions
			for name, param := range tool.Parameters {
				params["properties"].(map[string]interface{})[name] = map[string]interface{}{
					"type":        string(param.Type),
					"description": param.Description,
				}
				if param.Required {
					params["required"] = append(params["required"].([]string), name)
				}
			}

			tools[i] = lmStudioTool{
				Type: "function",
				Function: lmStudioToolFunction{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  params,
				},
			}
		}
	}

	// Build request
	lmReq := lmStudioRequest{
		Model:    l.Model,
		Messages: messages,
		Stream:   false,
	}

	// Only add tools if there are any
	if len(tools) > 0 {
		lmReq.Tools = tools
	}

	// Marshal request
	jsonBody, err := json.Marshal(lmReq)
	if err != nil {
		logging.Error("Failed to marshal LM Studio request", "error", err)
		return LLMResponse{}, err
	}

	// Log request details
	logging.Info("Sending request to LM Studio",
		"base_url", l.BaseURL,
		"model", l.Model,
		"system_prompt_length", len(request.SystemPrompt),
		"user_prompt_length", len(request.Prompt),
		"tools_count", len(tools),
	)
	logging.Debug("LM Studio request details",
		"system_prompt", request.SystemPrompt,
		"user_prompt", request.Prompt,
	)

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", l.BaseURL+"/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		logging.Error("Failed to create HTTP request", "error", err)
		return LLMResponse{}, err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+l.APIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		logging.Error("Failed to send request to LM Studio", "error", err)
		return LLMResponse{}, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read response body", "error", err)
		return LLMResponse{}, err
	}

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		logging.Error("LM Studio returned non-200 status",
			"status_code", resp.StatusCode,
			"body", string(body),
		)
		return LLMResponse{
			StatusCode: resp.StatusCode,
			Response:   fmt.Sprintf("LM Studio error: %s", string(body)),
		}, nil
	}

	// Parse response
	var lmResp lmStudioResponse
	if err := json.Unmarshal(body, &lmResp); err != nil {
		logging.Error("Failed to unmarshal LM Studio response", "error", err, "body", string(body))
		return LLMResponse{}, err
	}

	// Extract response content and tool calls
	if len(lmResp.Choices) == 0 {
		return LLMResponse{
			StatusCode: resp.StatusCode,
			Response:   "",
		}, nil
	}

	choice := lmResp.Choices[0]

	// Extract tool uses
	var toolUses []ToolUse
	for _, toolCall := range choice.Message.ToolCalls {
		// Parse arguments JSON
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
			logging.Error("Failed to parse tool call arguments",
				"error", err,
				"tool_name", toolCall.Function.Name,
				"arguments", toolCall.Function.Arguments,
			)
			continue
		}

		toolUses = append(toolUses, ToolUse{
			ToolName: toolCall.Function.Name,
			ToolArgs: args,
		})
	}

	// Log response details
	logging.Info("Received response from LM Studio",
		"status_code", resp.StatusCode,
		"content_length", len(choice.Message.Content),
		"tool_calls", len(toolUses),
		"finish_reason", choice.FinishReason,
		"total_tokens", lmResp.Usage.TotalTokens,
	)

	return LLMResponse{
		StatusCode: resp.StatusCode,
		Response:   choice.Message.Content,
		ToolUses:   toolUses,
	}, nil
}
