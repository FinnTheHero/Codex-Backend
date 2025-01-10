package server

import (
	admin_handler "Codex-Backend/api/server/handlers/admin"
	auth_handler "Codex-Backend/api/server/handlers/auth"
	client_handler "Codex-Backend/api/server/handlers/client"
	"Codex-Backend/api/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisteredRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",           // Local
			"https://codex-reader.vercel.app", // Remote TODO: change this to include url from env later.
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"X-Requested-With",
			"Authorization",
			"Accept",
			"Acces-Control-Allow-Origin",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
	}))

	client := r.Group("/")
	{
		client.Use(middleware.RateLimiter)
		client.GET("/all", client_handler.FindAllNovels)
		client.GET("/:novel", client_handler.FindNovel)
		client.GET("/:novel/all", client_handler.FindAllChapters)
		client.GET("/:novel/:chapter", client_handler.FindChapter)
		client.GET("/:novel/:chapter/next-previous", client_handler.FindPreviousAndNextChapters)
	}

	admin := r.Group("/admin")
	{
		admin.Use(middleware.RateLimiter)
		admin.POST("/novel", middleware.ValidateToken(), middleware.IsAdmin(), admin_handler.CreateNovel)
		admin.POST("/:novel/chapter", middleware.ValidateToken(), middleware.IsAdmin(), admin_handler.CreateChapter)
	}

	auth := r.Group("/auth")
	{
		auth.Use(middleware.RateLimiter)
		auth.GET("/validate", middleware.ValidateToken(), auth_handler.ValidateToken)
		auth.POST("/login", auth_handler.LoginUser)
		auth.POST("/logout", auth_handler.LogoutUser)
		auth.POST("/register", auth_handler.RegisterUser)
	}
}
