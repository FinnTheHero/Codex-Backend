package server

import (
	admin_handler "Codex-Backend/api/server/handlers/admin"
	auth_handler "Codex-Backend/api/server/handlers/auth"
	client_handler "Codex-Backend/api/server/handlers/client"
	"Codex-Backend/api/server/middleware"

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
		admin.POST("/novel", middleware.ValidateToken(), middleware.IsAdmin(), admin_handler.CreateNovel)
		admin.POST("/:novel/chapter", middleware.ValidateToken(), middleware.IsAdmin(), admin_handler.CreateChapter)
	}

	auth := r.Group("/auth")
	{
		auth.GET("/validate", middleware.ValidateToken(), auth_handler.ValidateToken)
		auth.GET("/login", auth_handler.LoginUser)
		auth.POST("/register", auth_handler.RegisterUser)
	}

	// Change this to a more specific CORS policy in production
	// admin.Use(cors.Default())
}
