package handler

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/server/middleware/token"
	"Codex-Backend/api/internal/service"
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

	err := service.RegisterUser(user, ctx)
	if e, ok := err.(*cmn.Error); ok {
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

	user, err := service.LoginUser(credentials, ctx)
	if e, ok := err.(*cmn.Error); ok {
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

	config := token.DefaultTokenConfig()

	tokens, err := token.GenerateTokenPair(user.ID, user.Email, config)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Error logging in the user: " + err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", tokens.AccessToken, int(config.AccessTTL.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(config.RefreshTTL.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"type":     user.Type,
		},
		"message": "Login successful",
	})
}

func LogoutUser(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Error logging out the user: " + err.Error(),
		})
		return
	}

	err = service.LogoutUser(tokenString)
	if e, ok := err.(*cmn.Error); ok {
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

	c.SetCookie("access_token", "", -1, "", "", true, true)
	c.SetCookie("refresh_token", "", -1, "", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
