package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	TOKEN	string
    ENDPOINT string

}

func LoadConfig() Config {
    // path, err := os.Executable()
    path, err := os.Getwd() // gets current working dir (where `go run` is called)
    if err != nil {
        log.Fatalf("Failed to get path: %v", err)
    }

    dir := filepath.Dir(path)
    //parentDir := filepath.Dir(dir) // one level above the executable

    dotenvPath := filepath.Join(dir, ".env")

    if err := godotenv.Load(dotenvPath); err != nil {
        log.Println("No .env file found at", dotenvPath, ", relying on environment variables")
    }

    return Config{
        TOKEN:    os.Getenv("TOKEN"),
        ENDPOINT: os.Getenv("ENDPOINT"),
    }
}
