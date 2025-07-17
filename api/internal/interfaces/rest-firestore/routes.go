package firestore_server

import (
	firestore_handlers "Codex-Backend/api/internal/interfaces/rest-firestore/handlers"
	firestore_middleware "Codex-Backend/api/internal/interfaces/rest-firestore/middleware"
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

	r.Use(middleware.RateLimiter)

	client := r.Group("/")
	{
		client.Use(middleware.VerifyUsersTablesExist())
		client.GET("/all", firestore_handlers.FindAllNovels)
		client.GET("/:novel", firestore_handlers.FindNovel)
		client.GET("/:novel/all", firestore_handlers.FindAllChapters)
		client.GET("/:novel/:chapter", firestore_handlers.FindChapter)
		// client.GET("/:novel/:chapter/next-previous", handlers.FindPreviousAndNextChapters)
	}

	admin := r.Group("/admin")
	{
		admin.Use(firestore_middleware.ValidateToken())
		admin.Use(firestore_middleware.IsAdmin())
		admin.POST("/novel", firestore_handlers.CreateNovel)
		admin.POST("/:novel/chapter", firestore_handlers.CreateChapter)
	}

	auth := r.Group("/auth")
	{
		auth.GET("/validate", firestore_middleware.ValidateToken(), firestore_handlers.ValidateToken)
		auth.POST("/login", firestore_handlers.LoginUser)
		auth.POST("/logout", firestore_handlers.LogoutUser)
		auth.POST("/register", firestore_handlers.RegisterUser)
	}
}
