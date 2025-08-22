package server

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/server/handler"
	"Codex-Backend/api/internal/server/middleware"
	"Codex-Backend/api/internal/server/middleware/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisteredRoutes(r *gin.Engine) {
	domain_url := cmn.GetEnvVariable("DOMAIN")
	if gin.Mode() == gin.DebugMode && domain_url == "" {
		domain_url = "*"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			domain_url,
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

	r.Use(middleware.RateLimiter())
	r.MaxMultipartMemory = 32 << 20 // 32 MB

	token.InitIMTokenCache()

	// Add mandatory token check
	r.Use(token.SetClaimsFromToken(), token.GlobalToken.AutoRefreshTokenMiddleware())

	client := r.Group("/")
	{
		client.GET("/all", handler.FindAllNovels)
		client.GET("/:novel", handler.FindNovel)
		client.GET("/:novel/all", handler.FindAllChapters)
		client.GET("/:novel/:chapter", handler.FindChapter)
		client.GET("/:novel/chapters", handler.GetPaginatedChapters)
	}

	manage := r.Group("/manage")
	{
		manage.Use(token.GlobalToken.LoadUser())

		// Create
		manage.POST("/novel", handler.CreateNovel)
		manage.POST("/:novel/chapter", handler.CreateChapter)
		manage.POST("/epub", handler.EPUBNovel)

		// Update
		manage.PUT("/:novel", handler.UpdateNovel)
		manage.PUT("/:novel/:chapter", handler.UpdateChapter)

		// Delete
		manage.DELETE("/:novel", handler.DeleteNovel)
		manage.DELETE("/:novel/:chapter", handler.DeleteChapter)
	}

	user := r.Group("/user")
	{
		user.GET("/validate", handler.ValidateToken)
		user.POST("/login", handler.LoginUser)
		user.POST("/logout", handler.LogoutUser)
		user.POST("/register", handler.RegisterUser)
		user.GET("/refresh", handler.RefreshToken)
	}
}
