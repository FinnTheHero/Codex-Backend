package firestore_server

import (
	firestore_handlers "Codex-Backend/api/internal/interfaces/rest/handlers"
	firestore_middleware "Codex-Backend/api/internal/interfaces/rest/middleware"

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
			"PUT",
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

	r.Use(firestore_middleware.RateLimiter())

	client := r.Group("/")
	{
		client.GET("/all", firestore_handlers.FindAllNovels)
		client.GET("/:novel", firestore_handlers.FindNovel)
		client.GET("/:novel/all", firestore_handlers.FindAllChapters)
		client.GET("/:novel/:chapter", firestore_handlers.FindChapter)
	}

	manage := r.Group("/manage")
	{
		manage.POST("/novel", firestore_middleware.ValidateToken(), firestore_handlers.CreateNovel)
		manage.POST("/:novel/chapter", firestore_middleware.ValidateToken(), firestore_handlers.CreateChapter)
		manage.PUT("/:novel", firestore_middleware.ValidateToken(), firestore_handlers.UpdateNovel)
		manage.PUT("/:novel/:chapter", firestore_middleware.ValidateToken(), firestore_handlers.UpdateChapter)

	}

	user := r.Group("/user")
	{
		user.GET("/validate", firestore_middleware.ValidateToken(), firestore_handlers.ValidateToken)
		user.POST("/login", firestore_handlers.LoginUser)
		user.POST("/logout", firestore_handlers.LogoutUser)
		user.POST("/register", firestore_handlers.RegisterUser)
	}
}
