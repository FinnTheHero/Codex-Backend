package server

import (
	admin_handler "Codex-Backend/api/server/handlers/admin"
	client_handler "Codex-Backend/api/server/handlers/client"

	"github.com/gin-gonic/gin"
)

func RegisteredRoutes(r *gin.Engine) {
	client := r.Group("/")
	{
		client.GET("/all", client_handler.FindAllNovels)
		client.GET("/:novel", client_handler.FindNovel)
		client.GET("/:novel/all", client_handler.FindAllChapters)
		client.GET("/:novel/:chapter", client_handler.FindChapter)
	}

	admin := r.Group("/admin")
	{
		admin.POST("/novel", admin_handler.CreateNovel)
		admin.POST("/:novel/chapter", admin_handler.CreateChapter)
	}

	// Change this to a more specific CORS policy in production
	// admin.Use(cors.Default())
}
