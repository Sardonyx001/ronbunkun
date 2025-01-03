package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP HTTPConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		HTTP: LoadHTTPConfig(),
	}
}
