package llm

type LLMProvider interface {
	Generate(request LLMRequest) (LLMResponse, error)
}

type LLMResponse struct {
	StatusCode int
	Response   string
	ToolUses   []ToolUse
}

type ToolUse struct {
	ToolName string
	ToolArgs map[string]interface{}
}

type LLMRequest struct {
	Prompt string
	Tools  []Tool
}

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}
