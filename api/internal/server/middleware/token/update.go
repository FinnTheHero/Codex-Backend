package token

import (
	"Codex-Backend/api/internal/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (mf *IMTokenCache) UpdateAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		defer ctx.Done()

		claims, exists := c.Get("claims")
		if !exists {
			c.Next()
			return
		}

		userClaims, ok := claims.(*domain.Claims)
		if !ok {
			c.Next()
			return
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Missing or invalid refresh token",
			})
			return
		}

		config := DefaultTokenConfig()

		newAccessToken, err := refreshAccessTokenFromString(refreshToken, userClaims.ID, domain.LookupUser{
			Cache:         mf.cache,
			CacheDuration: 1 * time.Hour,
		}, ctx)
		if err != nil {
			c.Next()
			return
		}

		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("access_token", newAccessToken, int(config.AccessTTL.Seconds()), "/", "", true, true)

		c.Header("X-Token-Refreshed", "true")

		c.Next()
	}
}
