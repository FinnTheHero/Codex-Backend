package token

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (mf *IMTokenCache) LoadUser() gin.HandlerFunc {
	return LookupUser(domain.LookupUser{
		Cache:         mf.cache,
		CacheDuration: 1 * time.Hour,
	})
}

func LookupUser(config domain.LookupUser) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Check cache first
		var user *domain.User
		cacheKey := fmt.Sprintf("user:%s", userClaims.ID)

		if config.Cache != nil {
			if cached, found := config.Cache.Get(cacheKey); found {
				if cachedUser, ok := cached.(*domain.User); ok {
					user = cachedUser
				}
			}
		}

		// Fetch user if not in cache
		if user == nil {
			user, err := service.GetUserByID(userClaims.ID, c.Request.Context())
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "User verification failed",
				})
				return
			}

			// Cache user if cache is available
			if config.Cache != nil {
				config.Cache.Set(cacheKey, user, config.CacheDuration)
			}
		}

		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}
