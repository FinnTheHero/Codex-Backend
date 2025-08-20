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

		// Create
		manage.POST("/novel", token.AuthenticateOnly(), handler.CreateNovel)
		manage.POST("/:novel/chapter", token.AuthenticateOnly(), handler.CreateChapter)
		manage.POST("/epub", token.AuthenticateOnly(), handler.EPUBNovel)

		// Update
		manage.PUT("/:novel", token.AuthenticateOnly(), handler.UpdateNovel)
		manage.PUT("/:novel/:chapter", token.AuthenticateOnly(), handler.UpdateChapter)

		// Delete
		manage.DELETE("/:novel", token.AuthenticateOnly(), handler.DeleteNovel)
		manage.DELETE("/:novel/:chapter", token.AuthenticateOnly(), handler.DeleteChapter)
	}

	user := r.Group("/user")
	{
		user.GET("/validate", token.GlobalToken.AuthenticateAndLoadUser(), handler.ValidateToken)
		user.POST("/login", handler.LoginUser)
		user.POST("/logout", handler.LogoutUser)
		user.POST("/register", handler.RegisterUser)
		user.GET("/refresh", token.AuthenticateOnly())
	}
}
