package config

import (
	"errors"
	"os"
)

type Config struct {
	Port   string
	APIKey string
}

func Load() (*Config, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, errors.New("API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:   port,
		APIKey: apiKey,
	}, nil
}
