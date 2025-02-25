package handlers

import (
	"Codex-Backend/api/internal/domain"
	auth_service "Codex-Backend/api/internal/usecases/auth"
	error_service "Codex-Backend/api/internal/usecases/error"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {

	var user domain.NewUser

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

	err := auth_service.RegisterUser(user)
	if err != nil {
		statusError := http.StatusInternalServerError

		if errors.Is(err, error_service.ErrEmailTaken) {
			statusError = http.StatusConflict
		}

		if errors.Is(err, error_service.ErrUserNotFound) {
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


func LoginUser(c *gin.Context) {

	var credentials domain.Credentials

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

	token, user, err := auth_service.LoginUser(credentials)
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
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
			"type":     user.Type,
		},
		"authorized": true,
	})
	return
}

func LogoutUser(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

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
	user, ok := result.(domain.User)
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
