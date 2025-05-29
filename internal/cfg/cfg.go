package cfg

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	ApiKey   string
	BaseUrl  string
}

func ReadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	apiKey := os.Getenv("CEREBRAS_API_KEY")
	if apiKey == "" {
		log.Fatal("CEREBRAS_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("CEREBRAS_BASE_URL")
	if baseURL == "" {
		log.Fatal("CEREBRAS_BASE_URL environment variable is not set")
	}

	return Config{
		Port:     port,
		ApiKey:   apiKey,
		BaseUrl:  baseURL,
	}
}

func NewConfig(port, apiKey, baseUrl string) Config {
	return Config{
		Port:     port,
		ApiKey:   apiKey,
		BaseUrl:  baseUrl,
	}
}
