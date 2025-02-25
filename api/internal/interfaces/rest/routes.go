package server

import (
	"Codex-Backend/api/internal/interfaces/rest/handlers"
	"Codex-Backend/api/internal/interfaces/rest/middleware"

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
			"DELETE",
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
		client.GET("/all", handlers.FindAllNovels)
		client.GET("/:novel", handlers.FindNovel)
		client.GET("/:novel/all", handlers.FindAllChapters)
		client.GET("/:novel/:chapter", handlers.FindChapter)
		client.GET("/:novel/:chapter/next-previous", handlers.FindPreviousAndNextChapters)
	}

	admin := r.Group("/admin")
	{
		admin.Use(middleware.RateLimiter)
		admin.POST("/novel", middleware.ValidateToken(), middleware.IsAdmin(), handlers.CreateNovel)
		admin.POST("/:novel/chapter", middleware.ValidateToken(), middleware.IsAdmin(), handlers.CreateChapter)
	}

	auth := r.Group("/auth")
	{
		auth.Use(middleware.RateLimiter)
		auth.GET("/validate", middleware.ValidateToken(), handlers.ValidateToken)
		auth.POST("/login", handlers.LoginUser)
		auth.POST("/logout", handlers.LogoutUser)
		auth.POST("/register", handlers.RegisterUser)
	}
}
