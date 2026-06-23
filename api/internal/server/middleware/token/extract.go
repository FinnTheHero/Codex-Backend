package token

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Get JWT token from cookie or Authorization header
func ExtractToken(token_name string, c *gin.Context) (string, error) {
	// Try cookie first
	if tokenString, err := c.Cookie(token_name); err == nil {
		return tokenString, nil
	}

	// Try Authorization header as fallback
	authHeader := c.GetHeader(token_name)
	if authHeader == "" {
		return "", fmt.Errorf("no authorization token provided")
	}

	// Handle "Bearer <token>" format
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return authHeader, nil
}
