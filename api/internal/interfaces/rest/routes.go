package firestore_server

import (
	cmn "Codex-Backend/api/internal/common"
	firestore_handlers "Codex-Backend/api/internal/interfaces/rest/handlers"
	firestore_middleware "Codex-Backend/api/internal/interfaces/rest/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisteredRoutes(r *gin.Engine) {
	domain := cmn.GetEnvVariable("DOMAIN")
	if gin.Mode() == gin.DebugMode && domain == "" {
		domain = "*"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			domain,
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
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
			"Set-Cookie",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
			"Set-Cookie",
		},
		AllowCredentials: true,
	}))

	r.Use(firestore_middleware.RateLimiter())
	r.MaxMultipartMemory = 32 << 20 // 32 MB

	client := r.Group("/")
	{
		client.GET("/all", firestore_handlers.FindAllNovels)
		client.GET("/:novel", firestore_handlers.FindNovel)
		client.GET("/:novel/all", firestore_handlers.FindAllChapters)
		client.GET("/:novel/:chapter", firestore_handlers.FindChapter)
		client.GET("/:novel/chapters", firestore_handlers.GetPaginatedChapters)
	}

	manage := r.Group("/manage")
	{
		manage.POST("/novel", firestore_middleware.ValidateToken(), firestore_handlers.CreateNovel)
		manage.POST("/:novel/chapter", firestore_middleware.ValidateToken(), firestore_handlers.CreateChapter)
		manage.POST("/epub", firestore_middleware.ValidateToken(), firestore_handlers.EPUBNovel)
		manage.PUT("/:novel", firestore_middleware.ValidateToken(), firestore_handlers.UpdateNovel)
		manage.PUT("/:novel/:chapter", firestore_middleware.ValidateToken(), firestore_handlers.UpdateChapter)
		manage.DELETE("/:novel", firestore_middleware.ValidateToken(), firestore_handlers.DeleteNovel)
		manage.DELETE("/:novel/:chapter", firestore_middleware.ValidateToken(), firestore_handlers.DeleteChapter)
	}

	user := r.Group("/user")
	{
		user.GET("/validate", firestore_middleware.ValidateToken(), firestore_handlers.ValidateToken)
		user.POST("/login", firestore_handlers.LoginUser)
		user.POST("/logout", firestore_handlers.LogoutUser)
		user.POST("/register", firestore_handlers.RegisterUser)
	}
}
