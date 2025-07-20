package firestore_handlers

import (
	"Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_services "Codex-Backend/api/internal/usecases-firestore/collections"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	user := domain.NewUser{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error reading data: " + err.Error(),
		})
		return
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Username, email and password are required",
		})
		return
	}

	err := firestore_services.RegisterUser(user, ctx)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Error registering the user: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error registering the user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}

func LoginUser(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	credentials := domain.Credentials{}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error reading data: " + err.Error(),
		})
		return
	}

	if credentials.Email == "" || credentials.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	token, user, err := firestore_services.LoginUser(credentials, ctx)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Error logging in the user: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error logging in the user: " + err.Error(),
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
}

func LogoutUser(c *gin.Context) {
	err := firestore_services.LogoutUser(c)
	if e, ok := err.(*common.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "Error logging out the user: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error logging out the user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func ValidateToken(c *gin.Context) {
	result, ok := c.Get("user")
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user, ok := result.(*domain.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user structure",
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
}
