package auth_handler

import (
	"Codex-Backend/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateToken(c *gin.Context) {
	// Get user from context
	result, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Cast user to User struct
	user, ok := result.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Type assertion failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
			"type":     user.Type,
		},
		"authenticated": true,
	})
	return
}
