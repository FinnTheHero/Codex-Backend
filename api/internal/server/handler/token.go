package handler

import (
	"Codex-Backend/api/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateToken(c *gin.Context) {
	result_claims, ok := c.Get("claims")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":          "User claims not found",
			"orignal_claims": result_claims,
		})
		return
	}

	claims, ok := result_claims.(*domain.Claims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid claims structure",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       claims.ID,
		"email":    claims.Email,
		"username": claims.Username,
		"type":     claims.Type,
	})
}
