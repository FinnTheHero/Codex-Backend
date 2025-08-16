package firestore_server

import (
	cmn "Codex-Backend/api/internal/common"

	"github.com/gin-gonic/gin"
)

func Server() {
	mode, err := cmn.GetEnvVariable("GIN_MODE")
	if err != nil {
		panic(err)
	}

	gin.SetMode(mode)

	r := gin.Default()

	RegisteredRoutes(r)

	r.Run()
}
