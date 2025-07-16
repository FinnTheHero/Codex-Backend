package firestore_middleware

import (
	"Codex-Backend/api/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user domain.User

		// Get user from context
		result, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Cast user to User struct
		user, ok = result.(domain.User)
		if !ok || user.Type != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized access",
			})
			return
		}

		c.Next()
	}
}
