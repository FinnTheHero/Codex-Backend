package main

import (
	cmn "Codex-Backend/api/common"
	firestore_server "Codex-Backend/api/internal/server"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func init() {
	cmn.LoadEnvVariables()

	mode := cmn.GetEnvVariable("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	firestore_server.Server()
}
