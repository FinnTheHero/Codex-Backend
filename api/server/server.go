package server

import (
	"github.com/gin-gonic/gin"
)

func Server() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	RegisteredRoutes(r)

	r.Run()
}
