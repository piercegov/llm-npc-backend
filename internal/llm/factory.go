package llm

import (
	"fmt"
	"strings"

	"github.com/piercegov/llm-npc-backend/internal/cfg"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

// NewProvider creates an LLM provider based on the configuration
func NewProvider(config cfg.Config) (LLMProvider, error) {
	provider := strings.ToLower(config.LLMProvider)

	switch provider {
	case "ollama":
		logging.Info("Creating Ollama provider", "model", config.OllamaModel)
		return NewOllama("11434"), nil

	case "lmstudio", "lm-studio":
		logging.Info("Creating LM Studio provider",
			"base_url", config.LMStudioBaseURL,
			"model", config.LMStudioModel,
		)
		return NewLMStudio(config.LMStudioBaseURL, config.LMStudioModel, config.LMStudioAPIKey), nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", config.LLMProvider)
	}
}
