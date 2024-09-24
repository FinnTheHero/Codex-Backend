package auth_handler

import (
	"Codex-Backend/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {

	var credentials models.Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	token, err := authService.LoginUser(credentials)
	if err != nil {
		statusError := http.StatusInternalServerError

		if err.Error() == "User not found" {
			statusError = http.StatusNotFound
		}

		if err.Error() == "Incorrect password" {
			statusError = http.StatusUnauthorized
		}

		c.JSON(statusError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("Authorization", token, 3600*24, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Loogin successful",
	})
	return

}
