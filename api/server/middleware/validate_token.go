package middleware

import (
	aws_services "Codex-Backend/api/aws/services"
	"Codex-Backend/api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			c.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(os.Getenv("JWT_SIGN_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			// Check token expiration
			if time.Now().Unix() > claims.ExpiresAt.Time.Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token expired",
				})
				return
			}

			// Find user
			result, err := aws_services.GetUser(claims.Email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			userDTO, ok := result.(models.UserDTO)
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Error casting user",
				})
				return
			}

			// Check user email
			if userDTO.Email != claims.Email {
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
