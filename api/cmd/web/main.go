package main

import (
	cmn "Codex-Backend/api/common"
	db_client "Codex-Backend/api/internal/database/client"
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
	connStr := cmn.GetEnvVariable("DATABASE_URL")
	client, err := db_client.GetClient(connStr)
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
