package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
	TOKEN	string
    // DBHost     string
    // DBPort     string
    // DBUser     string
    // DBPassword string
    // DBName     string
    // GRPCPort   string
}

func LoadConfig() Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on environment variables")
    }

    return Config{
        TOKEN:     os.Getenv("TOKEN"),
        // DBPort:     os.Getenv("DB_PORT"),
        // DBUser:     os.Getenv("DB_USER"),
        // DBPassword: os.Getenv("DB_PASSWORD"),
        // DBName:     os.Getenv("DB_NAME"),
        // GRPCPort:   os.Getenv("GRPC_PORT"),
    }
}
