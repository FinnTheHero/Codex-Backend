package middleware

import (
	"Codex-Backend/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Get user from context
		result, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Cast user to User struct
		user, ok = result.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Type assertion failed",
			})
			return
		}

		if user.Type != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized access",
			})
			return
		}

		c.Next()
		return
	}
}
