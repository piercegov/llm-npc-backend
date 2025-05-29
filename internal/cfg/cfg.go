package cfg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	ApiKey   string
	BaseUrl  string
	LogLevel string
}

func ReadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables and defaults")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("CEREBRAS_API_KEY")
	if apiKey == "" {
		log.Fatal("CEREBRAS_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("CEREBRAS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.cerebras.ai"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return Config{
		Port:     port,
		ApiKey:   apiKey,
		BaseUrl:  baseURL,
		LogLevel: logLevel,
	}
}

func NewConfig(port, apiKey, baseUrl, logLevel string) Config {
	return Config{
		Port:     port,
		ApiKey:   apiKey,
		BaseUrl:  baseUrl,
		LogLevel: logLevel,
	}
}
