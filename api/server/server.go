package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// Change this to a more specific CORS policy in production
	r.Use(cors.Default())

	RegisteredRoutes(r)

	r.Run()
}
