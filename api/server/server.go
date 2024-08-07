package server

import (
	admin_handler "Codex-Backend/api/server/handlers/admin"
	client_handler "Codex-Backend/api/server/handlers/client"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() {
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	client := r.Group("/novels")
	{
		client.GET("/all", client_handler.FindAllNovels)
		client.GET("/:novel", client_handler.FindNovel)
		client.GET("/:novel/all", client_handler.FindAllChapters)
		client.GET("/:novel/:chapter", client_handler.FindChapter)
	}

	client.Use(cors.Default())

	admin := r.Group("/admin")
	{
		admin.POST("/novel", admin_handler.CreateNovel)
		admin.POST("/:novel/chapter", admin_handler.CreateChapter)
	}

	admin.Use(cors.Default())

	r.Run()
}
