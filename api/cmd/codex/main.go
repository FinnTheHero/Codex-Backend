package main

import (
	"Codex-Backend/api/internal/config"
	server "Codex-Backend/api/internal/interfaces/rest"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func init(){
	if gin.Mode() == gin.DebugMode {
		err := config.LoadEnvVariables()
		if err != nil {
			log.Fatal(fmt.Sprintf("Error loading enviromental variables:" + err.Error()))
		}
	}
}

func main() {
	server.Server()
}
