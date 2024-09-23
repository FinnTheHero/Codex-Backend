package auth_handler

import (
	"Codex-Backend/api/models"
	auth_services "Codex-Backend/api/server/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

var authService = auth_services.NewAuthService()

func RegisterUser(c *gin.Context) {

	var user models.NewUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username, email and password are required",
		})
		c.Abort()
		return
	}

	err := authService.RegisterUser(user)
	if err != nil {
		statusError := http.StatusInternalServerError

		if err.Error() == "Email already in use" {
			statusError = http.StatusConflict
		}

		if err.Error() == "User not found" {
			statusError = http.StatusNotFound
		}

		c.JSON(statusError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
	c.Abort()
	return
}
