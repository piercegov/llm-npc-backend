package tools

import (
	"context"
	"fmt"

	"github.com/piercegov/llm-npc-backend/internal/llm"
)

// ToolResult represents the outcome of executing a tool
type ToolResult struct {
	Success bool
	Message string
	Data    map[string]interface{} // Additional data that might be needed by the game engine
}

// ToolHandler is a function that handles tool execution
// Game engines can register their own handlers for custom tools
type ToolHandler func(ctx context.Context, npcID string, args map[string]interface{}) (ToolResult, error)

// ToolRegistry manages tool definitions and their handlers
type ToolRegistry struct {
	tools    map[string]llm.Tool
	handlers map[string]ToolHandler
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:    make(map[string]llm.Tool),
		handlers: make(map[string]ToolHandler),
	}
}

// RegisterTool registers a tool with its handler
func (r *ToolRegistry) RegisterTool(tool llm.Tool, handler ToolHandler) error {
	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name)
	}
	
	r.tools[tool.Name] = tool
	r.handlers[tool.Name] = handler
	return nil
}

// GetTools returns all registered tools for LLM consumption
func (r *ToolRegistry) GetTools() []llm.Tool {
	tools := make([]llm.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ExecuteTool executes a tool by name
func (r *ToolRegistry) ExecuteTool(ctx context.Context, npcID string, toolUse llm.ToolUse) (ToolResult, error) {
	handler, exists := r.handlers[toolUse.ToolName]
	if !exists {
		return ToolResult{
			Success: false,
			Message: fmt.Sprintf("unknown tool: %s", toolUse.ToolName),
		}, fmt.Errorf("unknown tool: %s", toolUse.ToolName)
	}
	
	// Validate arguments match expected parameters
	tool, _ := r.tools[toolUse.ToolName]
	if err := validateArgs(tool, toolUse.ToolArgs); err != nil {
		return ToolResult{
			Success: false,
			Message: err.Error(),
		}, err
	}
	
	return handler(ctx, npcID, toolUse.ToolArgs)
}

// validateArgs validates that the provided arguments match the expected parameters
func validateArgs(tool llm.Tool, args map[string]interface{}) error {
	// Check required parameters
	for name, param := range tool.Parameters {
		if param.Required {
			if _, exists := args[name]; !exists {
				return fmt.Errorf("missing required parameter: %s", name)
			}
		}
	}
	
	// Check for unexpected parameters
	for argName := range args {
		if _, exists := tool.Parameters[argName]; !exists {
			return fmt.Errorf("unexpected parameter: %s", argName)
		}
	}
	
	// TODO: Add type validation based on param.Type
	
	return nil
}