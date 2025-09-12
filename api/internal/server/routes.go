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
	domains := cmn.GetDomains("DOMAIN")
	if gin.Mode() == gin.DebugMode && len(domains) == 0 {
		domains = []string{"*"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: domains,
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

	client := r.Group("/api/")
	{
		client.Use(token.SetClaimsFromToken(), token.GlobalToken.UpdateAccessToken(), token.GlobalToken.LoadUser())

		// Potentially add user public profile view here as well.

		client.GET("/all", handler.FindAllNovels)
		client.GET("/:novel", handler.FindNovel)
		client.GET("/:novel/all", handler.FindAllChapters)
		client.GET("/:novel/:chapter", handler.FindChapter)
		client.GET("/:novel/chapters", handler.GetPaginatedChapters)
	}

	manage := r.Group("/api/manage")
	{
		manage.Use(token.SetClaimsFromToken(), token.GlobalToken.UpdateAccessToken(), token.GlobalToken.LoadUser())

		// Create Novel/Chapters from epub file.
		manage.POST("/epub", handler.EPUBNovel)

		// Create
		manage.POST("/create/novel", handler.CreateNovel)
		manage.POST("/create/:novel/chapter", handler.CreateChapter)

		// Update
		manage.PUT("/update/:novel", handler.UpdateNovel)
		manage.PUT("/update/:novel/:chapter", handler.UpdateChapter)

		// Delete
		manage.DELETE("/delete/:novel", handler.DeleteNovel)
		manage.DELETE("/delete/:novel/:chapter", handler.DeleteChapter)
	}

	user := r.Group("/api/user")
	{
		user.POST("/login", handler.LoginUser)
		user.POST("/logout", handler.LogoutUser)
		user.POST("/register", handler.RegisterUser)
	}

	validate := r.Group("/api/validate")
	{
		validate.Use(token.SetClaimsFromToken(), token.GlobalToken.UpdateAccessToken(), token.GlobalToken.LoadUser())

		validate.GET("/", handler.ValidateToken)
	}
}
