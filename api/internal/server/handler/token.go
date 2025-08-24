package handler

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	token_middleware "Codex-Backend/api/internal/server/middleware/token"
	"Codex-Backend/api/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	defer ctx.Done()

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "No refresh token provided"})
		return
	}

	config := token_middleware.DefaultTokenConfig()
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.SigningKey), nil
	})

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid refresh token claims"})
		return
	}

	user, err := service.GetUserByID(claims.Subject, ctx)
	if e, ok := err.(*cmn.Error); ok {
		c.AbortWithStatusJSON(e.StatusCode(), gin.H{
			"error": "User not found: " + e.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "User not found: " + err.Error(),
		})
		return
	}

	// Generate new token pair
	tokens, err := token_middleware.GenerateTokenPair(user, config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Token generation failed"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", tokens.AccessToken, int(config.AccessTTL.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, int(config.RefreshTTL.Seconds()), "/", "", true, true)

	c.JSON(200, gin.H{
		"message":    "Tokens refreshed successfully",
		"expires_at": tokens.ExpiresAt,
		"expires_in": int(config.AccessTTL.Seconds()),
	})

}

func ValidateToken(c *gin.Context) {
	result_claims, ok := c.Get("claims")
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	claims, ok := result_claims.(*domain.Claims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user structure",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    claims.ID,
		"email": claims.Email,
	})
}
