package main

import (
	"Codex-Backend/api/internal/config"
	server "Codex-Backend/api/internal/interfaces/rest"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	if gin.Mode() == gin.DebugMode {
		err := config.LoadEnvVariables()
		if err != nil {
			log.Fatalf("Error loading env file: %s", err.Error())
		}
	}
}

func main() {
	server.Server()
}
