package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURL string
}

// Global AppConfig
var AppConfig *Config

func init() {
	LoadConfig()
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() {
	monogURL := os.Getenv("MONGO_URL")
	if monogURL == "" {
		log.Fatal("MONGO_URL environment variable is required")
	}

	AppConfig = &Config{
		MongoURL: monogURL,
	}
}
