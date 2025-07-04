package cfg

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

type Config struct {
	SocketPath      string
	HTTPPort        string
	ApiKey          string
	BaseUrl         string
	LogLevel        string
	OllamaModel     string
	OllamaBaseURL   string
	LLMProvider     string
	LMStudioBaseURL string
	LMStudioModel   string
	LMStudioAPIKey  string
}

func ReadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		logging.Warn("Warning: .env file not found, using environment variables and defaults")
	}

	socketPath := os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/tmp/llm-npc-backend.sock"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080"
	}

	apiKey := os.Getenv("CEREBRAS_API_KEY")
	if apiKey == "" {
		logging.Warn("CEREBRAS_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("CEREBRAS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.cerebras.ai"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "qwen3:1.7b"
	}

	ollamaBaseURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaBaseURL == "" {
		ollamaBaseURL = "http://10.0.0.85:11434"
	}

	llmProvider := os.Getenv("LLM_PROVIDER")
	if llmProvider == "" {
		llmProvider = "ollama" // Default to Ollama for backward compatibility
	}

	lmStudioBaseURL := os.Getenv("LMSTUDIO_BASE_URL")
	if lmStudioBaseURL == "" {
		lmStudioBaseURL = "http://localhost:1234"
	}

	lmStudioModel := os.Getenv("LMSTUDIO_MODEL")
	if lmStudioModel == "" {
		lmStudioModel = "model" // Default model identifier
	}

	lmStudioAPIKey := os.Getenv("LMSTUDIO_API_KEY")
	if lmStudioAPIKey == "" {
		lmStudioAPIKey = "lm-studio" // Default API key for LM Studio
	}

	return Config{
		SocketPath:      socketPath,
		HTTPPort:        httpPort,
		ApiKey:          apiKey,
		BaseUrl:         baseURL,
		LogLevel:        logLevel,
		OllamaModel:     ollamaModel,
		OllamaBaseURL:   ollamaBaseURL,
		LLMProvider:     llmProvider,
		LMStudioBaseURL: lmStudioBaseURL,
		LMStudioModel:   lmStudioModel,
		LMStudioAPIKey:  lmStudioAPIKey,
	}
}

func NewConfig(socketPath, httpPort, apiKey, baseUrl, logLevel, ollamaModel string) Config {
	return Config{
		SocketPath:      socketPath,
		HTTPPort:        httpPort,
		ApiKey:          apiKey,
		BaseUrl:         baseUrl,
		LogLevel:        logLevel,
		OllamaModel:     ollamaModel,
		OllamaBaseURL:   "http://10.0.0.85:11434", // Default Ollama base URL
		LLMProvider:     "ollama", // Default for backward compatibility
		LMStudioBaseURL: "http://localhost:1234",
		LMStudioModel:   "model",
		LMStudioAPIKey:  "lm-studio",
	}
}
