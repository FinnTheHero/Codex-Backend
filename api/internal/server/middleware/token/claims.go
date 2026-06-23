package token

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetClaimsFromToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := ExtractToken("access_token", c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid authorization token",
			})
			return
		}

		// Parse and validate JWT
		claims, err := ParseAndValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			return
		}

		c.Set("claims", claims)
	}
}
