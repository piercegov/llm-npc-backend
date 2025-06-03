package cfg

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/piercegov/llm-npc-backend/internal/logging"
)

type Config struct {
	SocketPath  string
	ApiKey      string
	BaseUrl     string
	LogLevel    string
	OllamaModel string
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

	return Config{
		SocketPath:  socketPath,
		ApiKey:      apiKey,
		BaseUrl:     baseURL,
		LogLevel:    logLevel,
		OllamaModel: ollamaModel,
	}
}

func NewConfig(socketPath, apiKey, baseUrl, logLevel, ollamaModel string) Config {
	return Config{
		SocketPath:  socketPath,
		ApiKey:      apiKey,
		BaseUrl:     baseUrl,
		LogLevel:    logLevel,
		OllamaModel: ollamaModel,
	}
}
