package tools

import (
	"context"
	"fmt"

	"github.com/piercegov/llm-npc-backend/internal/llm"
)

// ToolProvider is the interface for anything that can provide and execute tools
type ToolProvider interface {
	GetTools() []llm.Tool
	ExecuteTool(ctx context.Context, npcID string, toolUse llm.ToolUse) (ToolResult, error)
}

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

// NewToolRegistry creates a new tool registry with built-in tools
func NewToolRegistry() *ToolRegistry {
	registry := &ToolRegistry{
		tools:    make(map[string]llm.Tool),
		handlers: make(map[string]ToolHandler),
	}

	// Register built-in continue_thinking tool
	continueThinkingTool := llm.Tool{
		Name:        "continue_thinking",
		Description: "Signal that you want to continue thinking and processing after seeing tool results. Use this when you need to analyze tool outputs or make follow-up decisions. You will continue thinking immediately - the only new input will be results of tool calls from the current/previous iteration. Use this sparingly and only in cases when you need to wait for the results of another tool call.",
		Parameters: map[string]llm.ToolParameter{
			"reason": {
				Type:        llm.TypeString,
				Description: "Brief explanation of why you want to continue thinking",
				Required:    false,
			},
		},
	}

	// Handler that just returns success - the actual logic is in ActForTick
	continueThinkingHandler := func(ctx context.Context, npcID string, args map[string]interface{}) (ToolResult, error) {
		reason := "No reason provided"
		if r, ok := args["reason"].(string); ok {
			reason = r
		}
		return ToolResult{
			Success: true,
			Message: fmt.Sprintf("Continuing thinking: %s", reason),
		}, nil
	}

	registry.tools[continueThinkingTool.Name] = continueThinkingTool
	registry.handlers[continueThinkingTool.Name] = continueThinkingHandler

	return registry
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

// GetToolsWithSession returns all registered tools plus session-specific tools
func (r *ToolRegistry) GetToolsWithSession(sessionTools []llm.Tool) []llm.Tool {
	// Start with global tools
	tools := r.GetTools()

	// Add session tools
	tools = append(tools, sessionTools...)

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

// CombinedToolRegistry wraps a ToolRegistry with session-specific tools
type CombinedToolRegistry struct {
	base         *ToolRegistry
	sessionTools []llm.Tool
}

// NewCombinedToolRegistry creates a new combined registry
func NewCombinedToolRegistry(base *ToolRegistry, sessionTools []llm.Tool) *CombinedToolRegistry {
	return &CombinedToolRegistry{
		base:         base,
		sessionTools: sessionTools,
	}
}

// GetTools returns all tools (global + session)
func (c *CombinedToolRegistry) GetTools() []llm.Tool {
	return c.base.GetToolsWithSession(c.sessionTools)
}

// ExecuteTool tries to execute from session tools first, then falls back to global
func (c *CombinedToolRegistry) ExecuteTool(ctx context.Context, npcID string, toolUse llm.ToolUse) (ToolResult, error) {
	// For now, session tools are definition-only (executed by game engine)
	// So we just check if it exists in session tools
	for _, tool := range c.sessionTools {
		if tool.Name == toolUse.ToolName {
			// Return a result indicating the game engine should handle this
			return ToolResult{
				Success: true,
				Message: fmt.Sprintf("Session tool '%s' called - to be executed by game engine", toolUse.ToolName),
				Data: map[string]interface{}{
					"session_tool": true,
					"tool_name":    toolUse.ToolName,
					"args":         toolUse.ToolArgs,
				},
			}, nil
		}
	}

	// Fall back to global tools
	return c.base.ExecuteTool(ctx, npcID, toolUse)
}
