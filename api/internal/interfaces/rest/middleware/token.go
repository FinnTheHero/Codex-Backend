package firestore_middleware

import (
	"Codex-Backend/api/internal/common"
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure/collections"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

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

			key := cmn.GetEnvVariable("JWT_SIGN_KEY")
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

			client, err := firestore_client.FirestoreClient()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			defer client.Close()

			cl := firestore_collections.Client{Client: client}

			// Find user
			user, err := cl.GetUserByEmail(claims.Email, ctx)
			if e, ok := err.(*common.Error); ok {
				c.AbortWithStatusJSON(e.StatusCode(), gin.H{
					"error": "Error Verifying Cookie: " + e.Error(),
				})
				return
			}

			if user == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}

			if user.Email != claims.Email {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
				return
			}

			// Set user in context
			c.Set("user", user)

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
