
package main

import (
    "fmt"
    "my_project/backend/config"
)

func main() {
    cfg := config.LoadConfig()
    fmt.Println("DB running on", cfg.DBHost, "port", cfg.DBPort)
}

