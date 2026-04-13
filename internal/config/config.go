package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct{
	DatabaseURL string
	Port string
}

func Load() (*Config, error) {
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Println("Warning: .env not found")
    }

    config := &Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Port:        os.Getenv("PORT"),
    }
    log.Println("DATABASE_URL:", config.DatabaseURL)
    return config, nil
}