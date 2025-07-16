package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

type Config struct {
	TOKEN	string
    MODULE_NAME  string
    ENDPOINT string

}

func LoadConfig() Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on environment variables")
    }

    return Config{
        TOKEN:     os.Getenv("TOKEN"),
        MODULE_NAME:    os.Getenv("MODULE_NAME"),
        ENDPOINT: os.Getenv("ENDPOINT"),

    }
}
