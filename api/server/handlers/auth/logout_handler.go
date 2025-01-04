package auth_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutUser(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
