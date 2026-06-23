package main

import (
	cmn "Codex-Backend/api/common"
	db "Codex-Backend/api/internal/database"
	firestore_server "Codex-Backend/api/internal/server"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	cmn.LoadEnvVariables()

	mode := cmn.GetEnvVariable("GIN_MODE")
	gin.SetMode(mode)

	ctx := context.Background()
	client, err := db.GetClient(ctx)
	if err != nil {
		panic(fmt.Sprintf("db new client: %v", err))
	}
	if err := client.EnsureSchema(ctx); err != nil {
		panic(fmt.Sprintf("schema ensure failed: %v", err))
	}
}

func main() {
	firestore_server.Server()
}
