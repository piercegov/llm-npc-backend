package llm

// ParameterType represents the type of a tool parameter
type ParameterType string

const (
	TypeString  ParameterType = "string"
	TypeNumber  ParameterType = "number"
	TypeBoolean ParameterType = "boolean"
	TypeObject  ParameterType = "object"
	TypeArray   ParameterType = "array" // TODO: needs a secondary type for items
)

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
	SystemPrompt string
	Prompt       string
	Tools        []Tool
}

type Tool struct {
	Name        string
	Description string
	Parameters  map[string]ToolParameter
}

type ToolParameter struct {
	Type        ParameterType
	Description string
	Required    bool
}
