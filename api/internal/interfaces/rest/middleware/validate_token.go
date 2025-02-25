package middleware

import (
	"Codex-Backend/api/internal/config"
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			key, err := config.GetEnvVariable("JWT_SIGN_KEY")
			if err != nil {
				return nil, err
			}
			return []byte(key), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
			// Check token expiration
			if time.Now().After(time.Unix(claims.ExpiresAt.Unix(), 0)) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token expired",
				})
				return
			}

			if claims.Email == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Mail not found in token",
				})
			}

			// Find user
			result, err := repository.GetUser(claims.Email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			userDTO, ok := result.(domain.UserDTO)
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Error casting user",
				})
				return
			}

			// Check user email
			if userDTO.Email == "" || userDTO.Email != claims.Email {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}

			// Set user in context
			c.Set("user", userDTO.User)

			// Continue to handler
			c.Next()
			return
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
	}
}
