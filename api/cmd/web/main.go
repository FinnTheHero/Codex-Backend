package main

import (
	cmn "Codex-Backend/api/common"
	firestore_server "Codex-Backend/api/internal/server"
	"os"

	_ "github.com/heroku/x/hmetrics/onload"
)

func init() {
	if mode := os.Getenv("GIN_MODE"); mode == "debug" {
		cmn.LoadEnvVariables()
	}
}

func main() {
	firestore_server.Server()
}
