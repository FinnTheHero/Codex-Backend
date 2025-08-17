package firestore_server

import (
	cmn "Codex-Backend/api/internal/common"

	"github.com/gin-gonic/gin"
)

func Server() {
	mode := cmn.GetEnvVariable("GIN_MODE")

	gin.SetMode(mode)

	r := gin.Default()

	RegisteredRoutes(r)

	r.Run()
}
